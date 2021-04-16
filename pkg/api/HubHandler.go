package api

import (
	"encoding/json"
	"net/http"

	"github.com/prometheus/common/log"
)

type joinMessage struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Username  string `json:"username"`
	Anonymous bool   `json:"anonymous"`
}

// HubStatus Return status of all rooms
func (h BaseHandler) HubStatus(w http.ResponseWriter, r *http.Request) {
	resp := map[string]interface{}{}
	rooms := []map[string]interface{}{}
	users := []map[string]interface{}{}

	for _, v := range h.Hub.Rooms {
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
				"name": v.Username,
				"type": v.Type,
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
