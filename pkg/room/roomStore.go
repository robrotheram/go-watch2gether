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
	users := []*Meta{}
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

// func NewRoomStore(client *redis.Client) *RoomStore {
// 	return &RoomStore{client: client}
// }

// func (udb *RoomStore) Create(room Meta) error {
// 	ctx := context.Background()
// 	if _, err := udb.client.HSetNX(ctx, udb.GetRedisKey(room.ID), "id", &room).Result(); err != nil {
// 		log.Errorf("create: redis error: %w", err)
// 		return fmt.Errorf("create: redis error: %w", err)
// 	}
// 	return nil
// }

// func (udb *RoomStore) Find(roomID string, usePrefix bool) (*Meta, error) {
// 	ctx := context.Background()
// 	if usePrefix {
// 		roomID = udb.GetRedisKey(roomID)
// 	}
// 	result, err := udb.client.HGet(ctx, roomID, "id").Result()
// 	if err != nil && err != redis.Nil {
// 		return nil, fmt.Errorf("find: redis error: %w", err)
// 	}
// 	if result == "" {
// 		return nil, fmt.Errorf("find: not found")
// 	}

// 	room := &Meta{}
// 	if err := room.UnmarshalBinary([]byte(result)); err != nil {
// 		return nil, fmt.Errorf("find: unmarshal error: %w", err)
// 	}

// 	return room, nil
// }

// func (udb *RoomStore) FindByName(name string) (*Meta, error) {
// 	rooms, err := udb.GetAll()
// 	if err != nil {
// 		return nil, err
// 	}
// 	for _, r := range rooms {
// 		if r.Name == name {
// 			return r, nil
// 		}
// 	}
// 	return nil, fmt.Errorf("Room Not found")
// }

// func (udb *RoomStore) GetAll() ([]*Meta, error) {
// 	ctx := context.Background()
// 	iter := udb.client.Scan(ctx, 0, PREFIX+"*", 0).Iterator()
// 	rooms := []*Meta{}
// 	for iter.Next(ctx) {
// 		if room, err := udb.Find(iter.Val(), false); err == nil {
// 			rooms = append(rooms, room)
// 		}
// 	}
// 	return rooms, nil
// }
