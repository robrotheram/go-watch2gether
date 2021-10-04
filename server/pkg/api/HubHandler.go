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

type activeRoom struct {
	ID        string `json:"id"`
	BotActive bool   `json:"bot_active"`
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

	activeRooms := []activeRoom{}
	for _, r := range h.Hub.Rooms {
		activeRooms = append(activeRooms, activeRoom{
			ID:        r.ID,
			BotActive: r.Bot != nil,
		})
	}

	resp["active"] = activeRooms
	resp["rooms"] = rooms
	resp["users"] = users

	w.Header().Set("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(resp)
}
