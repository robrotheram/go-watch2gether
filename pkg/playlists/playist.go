package playlists

import (
	"log"
	"watch2gether/pkg/media"

	"github.com/asdine/storm"
	"github.com/google/uuid"
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
	*storm.DB
}

func NewPlaylistStore(store *storm.DB) *PlaylistStore {
	return &PlaylistStore{
		DB: store,
	}
}

func (store *PlaylistStore) GetAll() ([]Playlist, error) {
	var playlist []Playlist
	err := store.All(&playlist)
	return playlist, err
}

func (store *PlaylistStore) GetByUser(user string) ([]Playlist, error) {
	var playlist []Playlist
	err := store.Find("User", user, &playlist)
	return playlist, err
}

func (store *PlaylistStore) GetByChannel(channel string) ([]Playlist, error) {
	var playlist []Playlist
	err := store.Find("Channel", channel, &playlist)
	return playlist, err
}

func (store *PlaylistStore) GetById(id string) (Playlist, error) {
	var playlist Playlist
	err := store.One("ID", id, &playlist)
	return playlist, err
}

func (store *PlaylistStore) Create(playlist *Playlist) error {
	playlist.ID = uuid.NewString()
	playlist.RefeshList()
	return store.Save(playlist)
}

func (store *PlaylistStore) UpdatePlaylist(id string, playlist *Playlist) error {
	if _, err := store.GetById(id); err != nil {
		return err
	}
	playlist.RefeshList()
	return store.Update(playlist)
}

func (store *PlaylistStore) DeletePlaylist(id string) error {
	if playlist, err := store.GetById(id); err == nil {
		err := store.DeleteStruct(&playlist)
		return err
	} else {
		return err
	}
}
