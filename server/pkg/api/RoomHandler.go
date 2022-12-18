package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"watch2gether/pkg/room"
	meta "watch2gether/pkg/roomMeta"
	"watch2gether/pkg/user"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

func (h BaseHandler) GetRoomMeta(w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	id := vars["id"]
	room, err := h.Rooms.Find(id)
	if err != nil {
		return StatusError{http.StatusNotFound, fmt.Errorf("room does not exisit")}
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
		return StatusError{http.StatusBadRequest, fmt.Errorf("room does not exisit")}
	}

	var meta = meta.Meta{}
	err = json.NewDecoder(req.Body).Decode(&meta)
	if err != nil {
		return StatusError{http.StatusBadRequest, fmt.Errorf("unable to read message")}
	}
	//r.Update(meta)
	h.Rooms.Update(r)
	return nil
}

func (h BaseHandler) JoinRoom(w http.ResponseWriter, r *http.Request) error {
	usr, ok := r.Context().Value("user").(user.User)
	if !ok {
		return StatusError{http.StatusBadRequest, fmt.Errorf("unable to get user")}
	}

	var roomMsg = joinMessage{}
	err := json.NewDecoder(r.Body).Decode(&roomMsg)
	if err != nil {
		return StatusError{http.StatusBadRequest, fmt.Errorf("unable to read message")}
	}
	roomStr, _ := h.Rooms.GetOrCreate(roomMsg.ID, roomMsg.Name, usr)
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
		return StatusError{http.StatusBadRequest, fmt.Errorf("unable to read message")}
	}

	//Check user exists
	usr, err := h.Users.FindByName(roomMsg.Username)
	if err != nil {
		log.Error(err)
	}

	room, ok := h.Hub.FindRoom(roomMsg.ID)
	if !ok {
		return StatusError{http.StatusBadRequest, fmt.Errorf("room does not exisit")}
	}

	if !room.ContainsUserID(usr.ID) {
		return StatusError{http.StatusBadRequest, fmt.Errorf("user does not exisits")}
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
		return StatusError{http.StatusBadRequest, fmt.Errorf("room does not exisit")}
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
	meta, err := h.Datastore.Rooms.Find(id)
	if err != nil {
		return StatusError{http.StatusBadRequest, fmt.Errorf("room does not exisit")}
	}

	playistID := vars["id"]
	playlist, err := h.Playlist.Find(playistID)

	if err != nil {
		return StatusError{http.StatusBadRequest, fmt.Errorf("playlist does not exist")}
	}
	meta.Queue = append(meta.Queue, playlist.Videos...)
	h.Rooms.Update(meta)
	room, _ := h.Hub.FindRoom(id)
	room.Send(meta)
	return nil
}
