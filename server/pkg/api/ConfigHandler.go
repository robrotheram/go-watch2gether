package api

import (
	"encoding/json"
	"net/http"
)

type UIconfig struct {
	BotID string `json:"bot"`
}

func (h BaseHandler) GetConfig(w http.ResponseWriter, r *http.Request) error {
	config := UIconfig{
		BotID: h.Config.DiscordClientID,
	}
	w.Header().Set("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(config)
}
