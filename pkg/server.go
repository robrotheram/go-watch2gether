package pkg

import (
	"encoding/json"
	"fmt"
	"net/http"
	"watch2gether/pkg/room"
	client "watch2gether/pkg/room"
	user "watch2gether/pkg/user"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/prometheus/common/log"
)

type joinMessage struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Username string `json:"username"`
}

type BaseHandler struct {
	Hub   *Hub
	Users *user.UserStore
	Rooms *room.RoomStore
}

func (h BaseHandler) GetRoomMeta(w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	id := vars["id"]
	room, err := h.Rooms.Find(id)
	if err != nil {
		return StatusError{http.StatusNotFound, fmt.Errorf("Room Does not exisit")}
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(room)
	return nil
}

func (h BaseHandler) UpdateRoomMeta(w http.ResponseWriter, req *http.Request) error {
	vars := mux.Vars(req)
	w.Header().Set("Content-Type", "application/json")
	id := vars["id"]
	r, err := h.Rooms.Find(id)
	if err != nil {
		return StatusError{http.StatusBadRequest, fmt.Errorf("Room Does not exisit")}
	}

	var meta = room.Meta{}
	err = json.NewDecoder(req.Body).Decode(&meta)
	if err != nil {
		return StatusError{http.StatusBadRequest, fmt.Errorf("Unable to read message")}
	}
	r.Update(meta)
	h.Rooms.Update(r)
	return nil
}

func (h BaseHandler) JoinRoom(w http.ResponseWriter, r *http.Request) error {
	var roomMsg = joinMessage{}
	err := json.NewDecoder(r.Body).Decode(&roomMsg)
	if err != nil {
		return StatusError{http.StatusBadRequest, fmt.Errorf("Unable to read message")}
	}

	//Check user exists
	usr, err := h.Users.FindByField("Name", roomMsg.Username)
	if err != nil {
		usr = user.NewUser(roomMsg.Username)
		err := h.Users.Create(usr)
		if err != nil {
			log.Error(err)
		}

	}

	// Check if the room existis in REDIS. Create if not
	roomStr, err := h.Rooms.Find(roomMsg.Name)
	if err != nil {
		roomStr, err := h.Rooms.FindByField("Name", roomMsg.Name)
		if err != nil {
			log.Info("Room Not found. Making...")
			roomStr = room.NewMeta(roomMsg.Name, usr.ID)
			err := h.Rooms.Create(roomStr)
			if err != nil {
				log.Errorf("Room Create error:  %w", err)
				return err
			}
		}
	}

	//Check that this server is hosting the room
	hubRoom, ok := h.Hub.GetRoom(roomStr.ID)
	if !ok {
		hubRoom = room.New(roomStr, h.Rooms)
		h.Hub.AddRoom(hubRoom)
	}
	if hubRoom.Status != "Running" {
		h.Hub.StartRoom(roomStr.ID)
	}

	hubRoom.Join(usr)

	resp := map[string]interface{}{}
	resp["user"] = usr
	resp["room_id"] = hubRoom.ID

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)

	return nil
}

func (h BaseHandler) LeaveRoom(w http.ResponseWriter, r *http.Request) error {
	var roomMsg = joinMessage{}
	err := json.NewDecoder(r.Body).Decode(&roomMsg)
	if err != nil {
		return StatusError{http.StatusBadRequest, fmt.Errorf("Unable to read message")}
	}

	room, ok := h.Hub.FindRoom(roomMsg.Name)
	if !ok {
		return StatusError{http.StatusBadRequest, fmt.Errorf("Room does not exisit")}
	}

	if !room.ContainsUserID(roomMsg.ID) {
		return StatusError{http.StatusBadRequest, fmt.Errorf("User does not exisits")}
	}
	log.Info("USER LEFT")
	room.Leave(roomMsg.Username)
	return nil
}

func (h BaseHandler) DeleteRoom(w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	id := vars["id"]

	_, ok := h.Hub.FindRoom(id)
	if !ok {
		return StatusError{http.StatusBadRequest, fmt.Errorf("Room Does not exisit")}
	}
	h.Hub.DeleteRoom(id)
	return nil
}

func (h BaseHandler) ConnectRoom() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id := vars["id"]
		room, ok := h.Hub.FindRoom(id)
		if !ok {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("Error Room Unknown"))
			return
		}
		StartWebSocket(w, r, room)
	})
}

// HubStatus Return status of all rooms
func (h BaseHandler) HubStatus(w http.ResponseWriter, r *http.Request) {
	resp := map[string]interface{}{}
	rooms := []map[string]interface{}{}
	users := []map[string]interface{}{}

	for _, v := range h.Hub.rooms {
		meta, _ := h.Rooms.Find(v.ID)
		rooms = append(rooms, map[string]interface{}{
			"id":     v.ID,
			"status": v.Status,
			"meta":   meta,
		})
	}

	usrs, err := h.Users.GetAll()
	if err == nil {
		for _, v := range usrs {
			users = append(users, map[string]interface{}{
				"id":   v.ID,
				"name": v.Name,
			})
		}
	} else {
		log.Errorf("DB Find error %v", err)
	}

	resp["rooms"] = rooms
	resp["users"] = users

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func StartServer(connection string, hub *Hub, userDB *user.UserStore, rooms *room.RoomStore) error {
	api := BaseHandler{Hub: hub, Users: userDB, Rooms: rooms}

	r := mux.NewRouter()

	r.Handle("/api/v1/room/{id}/ws", api.ConnectRoom())
	r.Handle("/api/v1/room/{id}", Handler{api.GetRoomMeta}).Methods("GET")
	r.Handle("/api/v1/room/{id}", Handler{api.UpdateRoomMeta}).Methods("POST")
	r.Handle("/api/v1/room/{id}", Handler{api.DeleteRoom}).Methods("DELETE")

	r.Handle("/api/v1/room/{roomName}/join", Handler{api.JoinRoom}).Methods("POST")
	r.Handle("/api/v1/room/{roomName}/leave", Handler{api.LeaveRoom}).Methods("POST")

	r.Handle("/api/v1/scrape", Handler{getPageInfo}).Methods("GET")

	r.Handle("/api/v1/status/{roomName}", Handler{api.GetRoomMeta})
	r.HandleFunc("/api/v1/status", api.HubStatus)

	spa := spaHandler{staticPath: "ui/build", indexPath: "index.html"}
	r.PathPrefix("/").Handler(spa)

	headers := handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"})
	methods := handlers.AllowedMethods([]string{"GET", "POST", "DELETE", "PUT", "HEAD", "OPTIONS"})
	origins := handlers.AllowedOrigins([]string{"*"})

	return http.ListenAndServe(connection, handlers.CORS(headers, methods, origins)(r))
}

func getPageInfo(w http.ResponseWriter, r *http.Request) error {
	vars := r.URL.Query()
	url := vars["url"][0]

	s, err := Scrape(url, 1)
	if err != nil {
		return err
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(s.Preview)

	return nil
}

func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
}
func StartWebSocket(w http.ResponseWriter, req *http.Request, r *room.Room) error {
	enableCors(&w)
	vars := req.URL.Query()
	token := vars["token"][0]

	log.Info("TOKEN: " + token)

	socket, err := client.Upgrader.Upgrade(w, req, nil)
	if err != nil {
		return StatusError{http.StatusBadRequest, fmt.Errorf("Connection is not using the websocket protocol")}
	}
	client.NewClient(r, socket, token)
	return nil
}
