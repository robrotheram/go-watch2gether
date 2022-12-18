package playlist

import (
	"fmt"
	"watch2gether/pkg/media"

	"github.com/asdine/storm"
)

var PREFIX = "playlist"

type PlayistStore struct {
	session *storm.DB
}

type Playist struct {
	ID       string        `rethinkdb:"id,omitempty" json:"id"`
	Username string        `json:"username"`
	Name     string        `json:"name"`
	RoomID   string        `json:"room"`
	Videos   []media.Media `json:"videos"`
}

func NewPlayistStore(session *storm.DB) *PlayistStore {
	rs := &PlayistStore{session: session}
	return rs
}
func (udb *PlayistStore) Create(playlist *Playist) error {
	return udb.session.Save(playlist)
}

func (udb *PlayistStore) GetAll() ([]*Playist, error) {
	users := []*Playist{}
	err := udb.session.All(&users)
	return users, err
}

func (udb *PlayistStore) Find(id string) (*Playist, error) {
	var playist *Playist
	err := udb.session.Find("ID", id, playist)
	return playist, err
}

func (udb *PlayistStore) FindByRoomID(roomID string) ([]Playist, error) {
	playists, err := udb.GetAll()
	if err != nil {
		return []Playist{}, err
	}
	filteredPlaylist := []Playist{}
	for _, playist := range playists {
		if playist.RoomID == roomID {
			filteredPlaylist = append(filteredPlaylist, *playist)
		}
	}
	return filteredPlaylist, fmt.Errorf("Playlist not found")
}

func (udb *PlayistStore) Update(playlist *Playist) error {
	return udb.session.Save(playlist)
}

func (udb *PlayistStore) Delete(playist Playist) error {
	return udb.session.DeleteStruct(playist)
}
