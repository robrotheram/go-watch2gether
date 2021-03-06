package api

import (
	"fmt"
	"net/http"
	"watch2gether/pkg/room"

	"github.com/prometheus/common/log"
)

func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
}

func StartWebSocket(w http.ResponseWriter, req *http.Request, r *room.Room) error {
	enableCors(&w)
	vars := req.URL.Query()
	token := vars["token"][0]

	log.Info("TOKEN: " + token)

	socket, err := room.Upgrader.Upgrade(w, req, nil)
	if err != nil {
		return StatusError{http.StatusBadRequest, fmt.Errorf("connection is not using the websocket protocol")}
	}
	room.NewClient(r, socket, token)
	return nil
}
