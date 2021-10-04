package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"watch2gether/pkg/events"
	"watch2gether/pkg/media"
	"watch2gether/pkg/user"

	"github.com/gorilla/mux"
)

type VideoMsg struct {
	Url string `json:"url"`
}

func (h BaseHandler) AddVideo(w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	id := vars["id"]

	room, ok := h.Hub.FindRoom(id)
	if !ok {
		return StatusError{http.StatusBadRequest, fmt.Errorf("room not active")}
	}

	meta, err := h.Datastore.Rooms.Find(id)
	if err != nil {
		return StatusError{http.StatusBadRequest, fmt.Errorf("room does not exisit")}
	}

	var msg = VideoMsg{}
	err = json.NewDecoder(r.Body).Decode(&msg)
	if err != nil {
		return StatusError{http.StatusBadRequest, fmt.Errorf("unable to read message")}
	}
	usr := getUser(r)
	videos, err := media.NewVideo(msg.Url, usr.Username)
	if err != nil {
		return err
	}

	meta.Queue = append(meta.Queue, videos...)
	room.HandleEvent(events.Event{
		Action:  events.EVNT_UPDATE_QUEUE,
		Watcher: user.NewWatcher(usr),
		Queue:   meta.Queue,
	})

	return nil
}
