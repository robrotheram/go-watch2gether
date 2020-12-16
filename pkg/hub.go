package pkg

import (
	"encoding/json"
	"net/http"
	"time"

	log "github.com/sirupsen/logrus"
)

type hub struct {
	rooms map[string]*room
}

// Hub Global var storing all rooms.
// TODO: turn into store backed by DB
var Hub = hub{rooms: map[string]*room{}}

/*
DeleteRoom
Sends Delete event across all clients and then closes connection before deleting
*/
func (h *hub) DeleteRoom(roomName string) {
	log.Info("DELETING_ROOM")
	evt := Event{
		Action: "ROOM_EXIT",
	}
	b, _ := json.Marshal(evt)
	h.rooms[roomName].forward <- b
	h.rooms[roomName].quit <- true
	delete(h.rooms, roomName)
}

// CleanUP Every 5 seconds go through each room and check to see if there was a delete
func CleanUP() {
	log.Info("Staring Cleanup Routine")
	for {
		time.Sleep(5 * time.Second)
		log.Info("Checking Room Infomation")
		for _, room := range Hub.rooms {
			room.PurgeUsers()
			if len(room.Meta.Users) == 0 {
				Hub.DeleteRoom(room.Meta.Name)
			}
		}
	}
}

// HubStatus Return status of all rooms
func HubStatus(w http.ResponseWriter, r *http.Request) {
	resp := []map[string]interface{}{}
	for _, v := range Hub.rooms {
		resp = append(resp, map[string]interface{}{
			"id":     v.ID,
			"status": v.status,
			"meta":   v.Meta,
		})
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
