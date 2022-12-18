package meta

import (
	"fmt"
	user "watch2gether/pkg/user"

	"github.com/asdine/storm"
	log "github.com/sirupsen/logrus"
)

var PREFIX = "room"

type RoomStore struct {
	session *storm.DB
}

func (u *RoomStore) GetRedisKey(id string) string {
	return fmt.Sprintf("%s:%s", PREFIX, id)
}

func NewRoomStore(session *storm.DB) *RoomStore {
	rs := &RoomStore{session: session}
	rs.Cleanup()
	return rs
}

func (udb *RoomStore) Create(room *Meta) error {
	return udb.session.Save(room)
}

func (udb *RoomStore) GetAll() ([]*Meta, error) {
	rooms := []*Meta{}
	err := udb.session.All(&rooms)
	if err != nil {
		return rooms, err
	}
	return rooms, nil
}

func (udb *RoomStore) Find(id string) (*Meta, error) {
	return udb.FindByField("ID", id)
}

func (udb *RoomStore) FindByField(feild, value string) (*Meta, error) {
	var room []Meta
	err := udb.session.Find(feild, value, &room)
	if err != nil {
		return nil, err
	}
	if len(room) <= 0 {
		return nil, fmt.Errorf("room not found")
	}
	return &room[0], err
}

func (udb *RoomStore) Update(meta *Meta) error {
	if meta == nil {
		return fmt.Errorf("can not update meta if nil")
	}
	return udb.session.Save(meta)
}

func (udb *RoomStore) Delete(meta *Meta) error {
	return udb.session.DeleteStruct(meta)
}

func (udb *RoomStore) Cleanup() {
	rooms, err := udb.GetAll()
	if err != nil {
		return
	}
	for _, r := range rooms {
		if r.Owner == "" {
			udb.Delete(r)
		}
	}
}

func (rooms *RoomStore) GetOrCreate(roomID string, roomName string, usr user.User) (*Meta, error) {
	roomMeta, err := rooms.Find(roomID)
	if err != nil {
		roomMeta, err = rooms.FindByField("Name", roomName)
		if err != nil || roomMeta == nil {
			log.Info("Room Not found. Making...")
			roomMeta = NewMeta(roomName, usr)
			if roomID != "" {
				roomMeta.ID = roomID
			}
			err := rooms.Create(roomMeta)
			if err != nil {
				return roomMeta, fmt.Errorf("room Create error:  %w", err)
			}
		}
	}
	return roomMeta, nil
}
