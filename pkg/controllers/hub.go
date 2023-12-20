package controllers

import (
	"fmt"
	"w2g/pkg/playlists"

	"github.com/asdine/storm"
)

type Hub struct {
	channels  map[string]*Controller
	playlists *playlists.PlaylistStore
	db        *storm.DB
}

func NewHub(db *storm.DB) *Hub {
	return &Hub{
		channels:  make(map[string]*Controller),
		playlists: playlists.NewPlaylistStore(db),
		db:        db,
	}
}

func (hub *Hub) Add(id string, controller *Controller) {
	hub.channels[id] = controller
}

func (hub *Hub) New(id string) *Controller {
	hub.channels[id] = NewController(id, hub.db)
	return hub.channels[id]
}

func (hub *Hub) Get(id string) (*Controller, error) {
	if _, ok := hub.channels[id]; !ok {
		return nil, fmt.Errorf("channel not found")
	}
	return hub.channels[id], nil
}

func (hub *Hub) Playlists() *playlists.PlaylistStore {
	return hub.playlists
}
