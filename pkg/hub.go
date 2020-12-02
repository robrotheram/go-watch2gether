package pkg

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type hub struct {
	rooms map[string]*room
}

var Hub = hub{rooms: map[string]*room{}}

func (h *hub) DeleteRoom(roomName string) {
	fmt.Println("DELETING_ROOM")
	h.rooms[roomName].forward <- []byte("ROOM_EXIT")
	h.rooms[roomName].quit <- true
	delete(h.rooms, roomName)
}

func CleanUP() {
	fmt.Println("Staring Cleanup Routine")
	for {
		time.Sleep(5 * time.Second)
		fmt.Println("Staring Cleanup Routine")
		for _, room := range Hub.rooms {
			room.PurgeUsers()
			if len(room.Meta.Users) == 0 {
				Hub.DeleteRoom(room.Meta.Name)
			}
		}
	}
}

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
