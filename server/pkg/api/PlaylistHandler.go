package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"watch2gether/pkg/media"

	"github.com/gorilla/mux"
)

func (h BaseHandler) GetRoomPlaylist(w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	id := vars["id"]
	room, err := h.Playlist.Find(id)
	if err != nil {
		return StatusError{http.StatusNotFound, fmt.Errorf("Room Does not exisit")}
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
		return StatusError{http.StatusNotFound, fmt.Errorf("Playists Does not exisit: %v", err)}
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
		return StatusError{http.StatusBadRequest, fmt.Errorf("Playlist not found message")}
	}
	return h.Playlist.Delete(id)
}

func (h BaseHandler) CretePlaylist(w http.ResponseWriter, req *http.Request) error {
	vars := mux.Vars(req)
	w.Header().Set("Content-Type", "application/json")
	id := vars["id"]
	r, err := h.Rooms.Find(id)
	if err != nil {
		return StatusError{http.StatusBadRequest, fmt.Errorf("Room Does not exisit")}
	}

	var playlist = media.Playist{}
	err = json.NewDecoder(req.Body).Decode(&playlist)
	if err != nil {
		return StatusError{http.StatusBadRequest, fmt.Errorf("Unable to read message")}
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
		return StatusError{http.StatusBadRequest, fmt.Errorf("Room Does not exisit")}
	}

	var playlist = media.Playist{}
	err = json.NewDecoder(req.Body).Decode(&playlist)
	if err != nil {
		return StatusError{http.StatusBadRequest, fmt.Errorf("Unable to read message")}
	}
	if r.ID != playlist.RoomID {
		return StatusError{http.StatusBadRequest, fmt.Errorf("Rooms do not match")}
	}
	return h.Playlist.Update(&playlist)
}
