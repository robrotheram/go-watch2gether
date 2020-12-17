package pkg

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

type joinMessage struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Username string `json:"username"`
}

func GetRoomMeta(w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	roomName := vars["roomName"]
	if _, ok := Hub.rooms[roomName]; !ok {
		return StatusError{http.StatusNotFound, fmt.Errorf("Room Does not exisit")}
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(Hub.rooms[roomName].Meta)
	return nil
}

func UpdateRoomMeta(w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	w.Header().Set("Content-Type", "application/json")
	roomName := vars["roomName"]
	if _, ok := Hub.rooms[roomName]; !ok {
		return StatusError{http.StatusBadRequest, fmt.Errorf("Room Does not exisit")}
	}

	var meta = roomMeta{}
	err := json.NewDecoder(r.Body).Decode(&meta)
	if err != nil {
		return StatusError{http.StatusBadRequest, fmt.Errorf("Unable to read message")}
	}
	Hub.rooms[roomName].Meta.Update(meta)
	return nil
}

func JoinRoom(w http.ResponseWriter, r *http.Request) error {
	var roomMsg = joinMessage{}
	err := json.NewDecoder(r.Body).Decode(&roomMsg)
	if err != nil {
		return StatusError{http.StatusBadRequest, fmt.Errorf("Unable to read message")}
	}

	if _, ok := Hub.rooms[roomMsg.Name]; !ok {
		log.Info("Creaeting New room:" + roomMsg.Name)
		Hub.rooms[roomMsg.Name] = newRoom(roomMsg.Name)
		go Hub.rooms[roomMsg.Name].run()
	}

	if Hub.rooms[roomMsg.Name].ContainsUser(roomMsg.Username) {
		return StatusError{http.StatusBadRequest, fmt.Errorf("User already exisits")}
	}

	if Hub.rooms[roomMsg.Name].status != "Running" {
		log.Info("Room Starting Starting" + roomMsg.Name)
		go Hub.rooms[roomMsg.Name].run()
	}

	Hub.rooms[roomMsg.Name].Join(NewUser(roomMsg.Username))

	return nil
}

func LeaveRoom(w http.ResponseWriter, r *http.Request) error {
	var roomMsg = joinMessage{}
	err := json.NewDecoder(r.Body).Decode(&roomMsg)
	if err != nil {
		return StatusError{http.StatusBadRequest, fmt.Errorf("Unable to read message")}
	}

	if _, ok := Hub.rooms[roomMsg.Name]; !ok {
		return StatusError{http.StatusBadRequest, fmt.Errorf("Room does not exisit")}
	}

	if !Hub.rooms[roomMsg.Name].ContainsUser(roomMsg.Username) {
		return StatusError{http.StatusBadRequest, fmt.Errorf("User does not exisits")}
	}
	fmt.Print("USER LEFT")
	Hub.rooms[roomMsg.Name].Leave(roomMsg.Username)

	return nil
}

func DeleteRoom(w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	roomName := vars["roomName"]

	if _, ok := Hub.rooms[roomName]; !ok {
		return StatusError{http.StatusBadRequest, fmt.Errorf("Room Does not exisit")}
	}
	Hub.DeleteRoom(roomName)
	return nil
}

func ConnectRoom() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		roomName := vars["roomName"]
		if _, ok := Hub.rooms[roomName]; !ok {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("Error Room Unknown"))
			return
		}
		Hub.rooms[roomName].ServeHTTP(w, r)
	})
}

func StartServer(connection string) error {
	r := mux.NewRouter()

	r.Handle("/api/v1/room/{roomName}/ws", ConnectRoom())

	r.Handle("/api/v1/room/{roomName}/join", Handler{JoinRoom}).Methods("POST")
	r.Handle("/api/v1/room/{roomName}/leave", Handler{LeaveRoom}).Methods("POST")

	r.Handle("/api/v1/room/{roomName}", Handler{GetRoomMeta}).Methods("GET")
	r.Handle("/api/v1/room/{roomName}", Handler{UpdateRoomMeta}).Methods("POST")
	r.Handle("/api/v1/room/{roomName}", Handler{DeleteRoom}).Methods("DELETE")

	r.Handle("/api/v1/scrape", Handler{getPageInfo}).Methods("GET")

	r.Handle("/api/v1/status/{roomName}", Handler{GetRoomMeta})
	r.HandleFunc("/api/v1/status", HubStatus)

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
