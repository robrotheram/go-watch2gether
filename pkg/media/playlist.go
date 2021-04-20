package media

import (
	"fmt"

	"gopkg.in/rethinkdb/rethinkdb-go.v6"
)

var PREFIX = "playlist"

type PlayistStore struct {
	session *rethinkdb.Session
}

type Playist struct {
	ID       string  `rethinkdb:"id,omitempty" json:"id"`
	Username string  `json:"username"`
	Name     string  `json:"name"`
	RoomID   string  `json:"room"`
	Videos   []Video `json:"videos"`
}

func NewPlayistStore(session *rethinkdb.Session) *PlayistStore {
	rs := &PlayistStore{session: session}
	return rs
}
func (udb *PlayistStore) Create(playlist *Playist) error {
	_, err := rethinkdb.Table(PREFIX).Insert(playlist).RunWrite(udb.session)
	return err
}

func (udb *PlayistStore) GetAll() ([]*Playist, error) {
	users := []*Playist{}
	// Fetch all the items from the database
	res, err := rethinkdb.Table(PREFIX).Run(udb.session)
	if err != nil {
		return users, err
	}
	err = res.All(&users)
	if err != nil {
		return users, err
	}
	return users, nil
}

func (udb *PlayistStore) Find(id string) (*Playist, error) {
	res, err := rethinkdb.Table(PREFIX).Get(id).Run(udb.session)
	if err != nil {
		return nil, err
	}
	if res.IsNil() {
		return nil, fmt.Errorf("playist not found")
	}
	var playlist *Playist
	res.One(&playlist)
	res.Close()
	return playlist, nil
}

func (udb *PlayistStore) FindByField(feild, value string) ([]Playist, error) {
	res, err := rethinkdb.Table(PREFIX).Filter(rethinkdb.Row.Field(feild).Eq(value)).Run(udb.session)
	var playlist []Playist
	if err != nil {
		return nil, err
	}

	if res.IsNil() {
		return []Playist{}, nil
	}

	res.All(&playlist)
	res.Close()
	return playlist, nil
}

func (udb *PlayistStore) Update(playlist *Playist) error {
	_, err := rethinkdb.Table(PREFIX).Get(playlist.ID).Update(playlist).RunWrite(udb.session)
	return err
}

func (udb *PlayistStore) Delete(id string) error {
	_, err := rethinkdb.Table(PREFIX).Get(id).Delete().RunWrite(udb.session)
	return err
}
