package datastore

import (
	"log"
	"watch2gether/pkg/hub"
	"watch2gether/pkg/media"
	meta "watch2gether/pkg/roomMeta"
	"watch2gether/pkg/user"
	"watch2gether/pkg/utils"
)

type Datastore struct {
	Rooms    *meta.RoomStore
	Playlist *media.PlayistStore
	Users    *user.UserStore
	Hub      *hub.Hub
}

func NewDatastore(config utils.Config) *Datastore {
	rethink, err := createSession(config)
	if err != nil {
		log.Fatalf("DB Error: %v", err)
	}

	hub := hub.NewHub()
	userStore := user.NewUserStore(rethink)
	createTable(rethink, config, user.PREFIX)

	roomStore := meta.NewRoomStore(rethink)
	createTable(rethink, config, meta.PREFIX)

	playlistStore := media.NewPlayistStore(rethink)
	createTable(rethink, config, media.PREFIX)

	return &Datastore{
		Rooms:    roomStore,
		Playlist: playlistStore,
		Users:    userStore,
		Hub:      hub,
	}
}

func (datastore *Datastore) StartCleanUP() {
	go datastore.Users.Cleanup()
	go datastore.Hub.CleanUP(datastore.Users)
}
