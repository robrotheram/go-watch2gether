package pkg

import (
	"encoding/json"
	"fmt"
	"net/http"

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
	Hub *Hub
}

func (h BaseHandler) GetRoomMeta(w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	roomName := vars["roomName"]
	room, ok := h.Hub.FindRoom(roomName)
	if !ok {
		return StatusError{http.StatusNotFound, fmt.Errorf("Room Does not exisit")}
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(room.Meta)
	return nil
}

func (h BaseHandler) UpdateRoomMeta(w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	w.Header().Set("Content-Type", "application/json")
	roomName := vars["roomName"]
	room, ok := h.Hub.FindRoom(roomName)
	if !ok {
		return StatusError{http.StatusBadRequest, fmt.Errorf("Room Does not exisit")}
	}

	var meta = roomMeta{}
	err := json.NewDecoder(r.Body).Decode(&meta)
	if err != nil {
		return StatusError{http.StatusBadRequest, fmt.Errorf("Unable to read message")}
	}
	room.Meta.Update(meta)
	return nil
}

func (h BaseHandler) JoinRoom(w http.ResponseWriter, r *http.Request) error {
	var roomMsg = joinMessage{}
	err := json.NewDecoder(r.Body).Decode(&roomMsg)
	if err != nil {
		return StatusError{http.StatusBadRequest, fmt.Errorf("Unable to read message")}
	}
	room, ok := h.Hub.FindRoom(roomMsg.Name)
	if !ok {
		log.Info("Room Not found. Making")
		room = h.Hub.NewRoom(roomMsg.Name)
	}
	if room.ContainsUser(roomMsg.Username) {
		return StatusError{http.StatusBadRequest, fmt.Errorf("User already exisits")}
	}
	if room.Status != "Running" {
		h.Hub.StartRoom(room.ID)
	}
	room.Join(NewUser(roomMsg.Username))
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

	if !room.ContainsUser(roomMsg.Username) {
		return StatusError{http.StatusBadRequest, fmt.Errorf("User does not exisits")}
	}
	log.Info("USER LEFT")
	room.Leave(roomMsg.Username)
	return nil
}

func (h BaseHandler) DeleteRoom(w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	roomName := vars["roomName"]

	_, ok := h.Hub.FindRoom(roomName)
	if !ok {
		return StatusError{http.StatusBadRequest, fmt.Errorf("Room Does not exisit")}
	}
	h.Hub.DeleteRoom(roomName)
	return nil
}

func (h BaseHandler) ConnectRoom() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		roomName := vars["roomName"]
		room, ok := h.Hub.FindRoom(roomName)
		if !ok {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("Error Room Unknown"))
			return
		}
		room.ServeHTTP(w, r)
	})
}

// HubStatus Return status of all rooms
func (h BaseHandler) HubStatus(w http.ResponseWriter, r *http.Request) {
	resp := []map[string]interface{}{}
	for _, v := range h.Hub.rooms {
		resp = append(resp, map[string]interface{}{
			"id":     v.ID,
			"status": v.Status,
			"meta":   v.Meta,
		})
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func StartServer(connection string, hub *Hub) error {
	api := BaseHandler{Hub: hub}

	r := mux.NewRouter()

	r.Handle("/api/v1/room/{roomName}/ws", api.ConnectRoom())

	r.Handle("/api/v1/room/{roomName}/join", Handler{api.JoinRoom}).Methods("POST")
	r.Handle("/api/v1/room/{roomName}/leave", Handler{api.LeaveRoom}).Methods("POST")

	r.Handle("/api/v1/room/{roomName}", Handler{api.GetRoomMeta}).Methods("GET")
	r.Handle("/api/v1/room/{roomName}", Handler{api.UpdateRoomMeta}).Methods("POST")
	r.Handle("/api/v1/room/{roomName}", Handler{api.DeleteRoom}).Methods("DELETE")

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
