package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	playlist "watch2gether/pkg/playlists"

	"github.com/gorilla/mux"
)

func (h BaseHandler) GetRoomPlaylist(w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	id := vars["id"]
	room, err := h.Playlist.Find(id)
	if err != nil {
		return StatusError{http.StatusNotFound, fmt.Errorf("room Does not exisit")}
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(room)
	return nil
}

func (h BaseHandler) GetAllRoomPlaylists(w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	id := vars["id"]
	playist, err := h.Playlist.FindByField("RoomID", id)
	if err != nil {
		return StatusError{http.StatusNotFound, fmt.Errorf("playists does not exisit: %v", err)}
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(playist)
	return nil
}

func (h BaseHandler) DeletePlaylist(w http.ResponseWriter, req *http.Request) error {
	vars := mux.Vars(req)
	w.Header().Set("Content-Type", "application/json")
	id := vars["id"]

	_, err := h.Playlist.Find(id)
	if err != nil {
		return StatusError{http.StatusBadRequest, fmt.Errorf("playlist not found message")}
	}
	return h.Playlist.Delete(id)
}

func (h BaseHandler) CretePlaylist(w http.ResponseWriter, req *http.Request) error {
	vars := mux.Vars(req)
	w.Header().Set("Content-Type", "application/json")
	id := vars["id"]
	r, err := h.Rooms.Find(id)
	if err != nil {
		return StatusError{http.StatusBadRequest, fmt.Errorf("room does not exisit")}
	}

	var playlist = playlist.Playist{}
	err = json.NewDecoder(req.Body).Decode(&playlist)
	if err != nil {
		return StatusError{http.StatusBadRequest, fmt.Errorf("unable to read message")}
	}
	playlist.RoomID = r.ID
	return h.Playlist.Create(&playlist)
}

func (h BaseHandler) UpdatePlaylist(w http.ResponseWriter, req *http.Request) error {
	vars := mux.Vars(req)
	w.Header().Set("Content-Type", "application/json")
	id := vars["room_id"]
	r, err := h.Rooms.Find(id)
	if err != nil {
		return StatusError{http.StatusBadRequest, fmt.Errorf("room does not exisit")}
	}

	var playlist = playlist.Playist{}
	err = json.NewDecoder(req.Body).Decode(&playlist)
	if err != nil {
		return StatusError{http.StatusBadRequest, fmt.Errorf("unable to read message")}
	}
	if r.ID != playlist.RoomID {
		return StatusError{http.StatusBadRequest, fmt.Errorf("rooms do not match")}
	}
	return h.Playlist.Update(&playlist)
}
