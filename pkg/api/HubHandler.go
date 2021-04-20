package api

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type joinMessage struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Username  string `json:"username"`
	Anonymous bool   `json:"anonymous"`
}

// HubStatus Return status of all rooms

func (h BaseHandler) HubStatus(w http.ResponseWriter, r *http.Request) error {
	resp := map[string]interface{}{}
	users := []map[string]interface{}{}

	rooms, _ := h.Rooms.GetAll()
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
		return fmt.Errorf("DB Find error %v", err)
	}

	resp["rooms"] = rooms
	resp["users"] = users

	w.Header().Set("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(resp)
}
