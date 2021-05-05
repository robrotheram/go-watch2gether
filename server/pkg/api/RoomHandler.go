package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"watch2gether/pkg/room"
	"watch2gether/pkg/user"

	"github.com/gorilla/mux"
	"github.com/prometheus/common/log"
)

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
	//r.Update(meta)
	h.Rooms.Update(r)
	return nil
}

func (h BaseHandler) JoinRoom(w http.ResponseWriter, r *http.Request) error {
	usr, ok := r.Context().Value("user").(user.User)
	if !ok {
		return StatusError{http.StatusBadRequest, fmt.Errorf("Unable to get user")}
	}

	var roomMsg = joinMessage{}
	err := json.NewDecoder(r.Body).Decode(&roomMsg)
	if err != nil {
		return StatusError{http.StatusBadRequest, fmt.Errorf("Unable to read message")}
	}
	roomStr, err := h.Rooms.GetOrCreate(roomMsg.ID, roomMsg.Name, usr)
	hubRoom, ok := h.Hub.GetRoom(roomStr.ID)
	if !ok {
		hubRoom = room.New(roomStr, h.Rooms)
		hubRoom.PurgeUsers(true)
		h.Hub.AddRoom(hubRoom)
	}

	found := hubRoom.ContainsUserID(usr.ID)

	if found {
		hubRoom.Disconnect(usr.ID)
		//hubRoom.Leave(usr.ID)
	}

	hubRoom.Join(usr)

	if hubRoom.Status != "Running" {
		h.Hub.StartRoom(roomStr.ID)
	}

	resp := map[string]interface{}{}
	resp["user"] = usr
	resp["room_id"] = hubRoom.ID

	w.Header().Set("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(resp)
}

func (h BaseHandler) LeaveRoom(w http.ResponseWriter, r *http.Request) error {
	var roomMsg = joinMessage{}
	err := json.NewDecoder(r.Body).Decode(&roomMsg)
	if err != nil {
		return StatusError{http.StatusBadRequest, fmt.Errorf("Unable to read message")}
	}

	//Check user exists
	usr, err := h.Users.FindByField("Name", roomMsg.Username)
	if err != nil {

		log.Error(err)

	}

	room, ok := h.Hub.FindRoom(roomMsg.ID)
	if !ok {
		return StatusError{http.StatusBadRequest, fmt.Errorf("Room does not exisit")}
	}

	if !room.ContainsUserID(usr.ID) {
		return StatusError{http.StatusBadRequest, fmt.Errorf("User does not exisits")}
	}

	log.Info("USER LEFT")
	room.Disconnect(usr.ID)
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

func (h BaseHandler) LoadPlaylist(w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	id := vars["room_id"]
	room, ok := h.Hub.FindRoom(id)
	if !ok {
		return StatusError{http.StatusBadRequest, fmt.Errorf("Room Does not exisit")}
	}
	playistID := vars["id"]
	playlist, err := h.Playlist.Find(playistID)

	if err != nil {
		return StatusError{http.StatusBadRequest, fmt.Errorf("Playlist does not exist")}
	}
	queue := room.GetQueue()
	queue = append(queue, playlist.Videos...)
	room.SetQueue(queue, user.SERVER_USER)
	return nil
}
