package datastore

import (
	"log"
	"os"
	"path/filepath"
	"watch2gether/pkg/hub"
	playlist "watch2gether/pkg/playlists"
	meta "watch2gether/pkg/roomMeta"
	"watch2gether/pkg/user"
	"watch2gether/pkg/utils"

	"github.com/asdine/storm"
)

type Datastore struct {
	Rooms    *meta.RoomStore
	Playlist *playlist.PlayistStore
	Users    *user.UserStore
	Hub      *hub.Hub
}

func NewDatastore(config utils.Config) *Datastore {

	os.MkdirAll("/app", os.ModePerm)
	path := filepath.Join(config.DatabasePath, "watch2gether.db")
	db, err := storm.Open(path)
	if err != nil {
		log.Fatalf("Unable to open db at: %s \n Error: %v", path, err)
	}

	hub := hub.NewHub()
	userStore := user.NewUserStore(db)
	roomStore := meta.NewRoomStore(db)
	playlistStore := playlist.NewPlayistStore(db)

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
