package room

import (
	"fmt"

	"gopkg.in/rethinkdb/rethinkdb-go.v6"
)

var PREFIX = "room"

type RoomStore struct {
	session *rethinkdb.Session
}

func (u *RoomStore) GetRedisKey(id string) string {
	return fmt.Sprintf("%s:%s", PREFIX, id)
}

func NewRoomStore(session *rethinkdb.Session) *RoomStore {
	rs := &RoomStore{session: session}
	rs.Cleanup()
	return rs
}

func (udb *RoomStore) Create(room *Meta) error {
	_, err := rethinkdb.Table(PREFIX).Insert(room).RunWrite(udb.session)
	return err
}

func (udb *RoomStore) GetAll() ([]*Meta, error) {
	rooms := []*Meta{}
	// Fetch all the items from the database
	res, err := rethinkdb.Table(PREFIX).Run(udb.session)
	if err != nil {
		return rooms, err
	}
	err = res.All(&rooms)
	if err != nil {
		return rooms, err
	}
	return rooms, nil
}

func (udb *RoomStore) Find(id string) (*Meta, error) {
	res, err := rethinkdb.Table(PREFIX).Get(id).Run(udb.session)
	if err != nil {
		return nil, err
	}
	if res.IsNil() {
		return nil, fmt.Errorf("User not found")
	}
	var room *Meta
	res.One(&room)
	res.Close()
	return room, nil
}

func (udb *RoomStore) FindByField(feild, value string) (*Meta, error) {
	res, err := rethinkdb.Table(PREFIX).Filter(rethinkdb.Row.Field(feild).Eq(value)).Run(udb.session)
	var room Meta
	if err != nil {
		return nil, err
	}

	if res.IsNil() {
		return nil, fmt.Errorf("room not found")
	}

	res.One(&room)
	res.Close()
	return &room, nil
}

func (udb *RoomStore) Update(meta *Meta) error {
	_, err := rethinkdb.Table(PREFIX).Get(meta.ID).Update(meta).RunWrite(udb.session)
	return err
}

func (udb *RoomStore) Delete(id string) error {
	_, err := rethinkdb.Table(PREFIX).Get(id).Delete().RunWrite(udb.session)
	return err
}

func (udb *RoomStore) Cleanup() {
	rooms, err := udb.GetAll()
	if err != nil {
		return
	}
	for _, r := range rooms {
		if r.Owner == "" {
			udb.Delete(r.ID)
		}
	}
}
