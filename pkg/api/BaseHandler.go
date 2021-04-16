package api

import (
	"net/http"
	"watch2gether/pkg/hub"
	"watch2gether/pkg/media"
	"watch2gether/pkg/room"
	"watch2gether/pkg/user"
)

type BaseHandler struct {
	Hub      *hub.Hub
	Users    *user.UserStore
	Rooms    *room.RoomStore
	Playlist *media.PlayistStore
}

// The Handler struct that takes a configured Env and a function matching
// our useful signature.
type Handler struct {
	H func(w http.ResponseWriter, r *http.Request) error
}
