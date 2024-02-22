package playlists

import (
	"log"
	"w2g/pkg/media"
	"w2g/pkg/utils"

	"github.com/google/uuid"
	bolt "go.etcd.io/bbolt"
)

type Playlist struct {
	ID      string        `json:"id"`
	Name    string        `json:"name"`
	Videos  []media.Media `json:"videos"`
	User    string        `json:"user"`
	Channel string        `json:"channel"`
}

func (playist *Playlist) RefeshList() {
	newList := []media.Media{}
	for _, video := range playist.Videos {
		if len(video.Type) > 0 {
			newList = append(newList, video)
			continue
		}
		if videos, err := media.NewVideo(video.Url, playist.User); err == nil {
			newList = append(newList, videos...)
		} else {
			log.Println(err)
		}
	}
	playist.Videos = newList
}

type PlaylistStore struct {
	*utils.Store[*Playlist]
}

func NewPlaylistStore(db *bolt.DB) *PlaylistStore {
	store := &utils.Store[*Playlist]{
		DB:     db,
		Bucket: []byte("controllers"),
	}
	store.Create()

	return &PlaylistStore{
		Store: store,
	}
}

func (store *PlaylistStore) GetAll() []*Playlist {
	return store.All()
}

func (store *PlaylistStore) GetByUser(user string) ([]*Playlist, error) {
	return store.Find("User", user)
}

func (store *PlaylistStore) GetByChannel(channel string) ([]*Playlist, error) {
	return store.Find("Channel", channel)
}

func (store *PlaylistStore) GetById(id string) (*Playlist, error) {
	return store.Get(id)
}

func (store *PlaylistStore) Create(playlist *Playlist) error {
	playlist.ID = uuid.NewString()
	playlist.RefeshList()
	return store.Save(playlist.ID, playlist)
}

func (store *PlaylistStore) UpdatePlaylist(id string, playlist *Playlist) error {
	if _, err := store.GetById(id); err != nil {
		return err
	}
	playlist.RefeshList()
	return store.Save(playlist.ID, playlist)
}

func (store *PlaylistStore) DeletePlaylist(id string) error {
	if playlist, err := store.GetById(id); err == nil {
		return store.Delete(playlist.ID)
	} else {
		return err
	}
}
