package user

import (
	"encoding/json"
	"fmt"

	"github.com/segmentio/ksuid"
	log "github.com/sirupsen/logrus"
	rethinkdb "gopkg.in/rethinkdb/rethinkdb-go.v6"
)

type User struct {
	ID   string `rethinkdb:"id,omitempty" json:"id"`
	Name string `json:"name"`
	Type string `json:"type"`
}

func (t *User) MarshalBinary() ([]byte, error) {
	return json.Marshal(t)
}

func (t *User) UnmarshalBinary(data []byte) error {
	if err := json.Unmarshal(data, &t); err != nil {
		return err
	}
	return nil
}

// func (u *User) AddRoom(id string) {
// 	u.Rooms = append(u.Rooms, id)
// }

// func (u *User) ConstainsRoom(id string) bool {
// 	for _, u := range u.Rooms {
// 		if u == id {
// 			return true
// 		}
// 	}
// 	return false
// }

// func (u *User) RemoveRoom(id string) {
// 	for i := range u.Rooms {
// 		room := u.Rooms[i]
// 		if room == id {
// 			u.Rooms = append(u.Rooms[:i], u.Rooms[i+1:]...)
// 			break
// 		}
// 	}
// }

func NewUser(name string) User {
	id := ksuid.New().String()
	log.Infof("Creating User with ID %s, Name: %s", id, name)
	return User{
		ID:   id,
		Name: name,
	}
}

var PREFIX = "user"

type UserStore struct {
	session *rethinkdb.Session
}

func (u *UserStore) GetKey(id string) string {
	return fmt.Sprintf("%s:%s", PREFIX, id)
}

func NewUserStore(session *rethinkdb.Session) *UserStore {
	return &UserStore{session: session}
}

func (udb *UserStore) Create(user User) error {
	_, err := rethinkdb.Table(PREFIX).Insert(user).RunWrite(udb.session)
	return err

}

func (udb *UserStore) GetAll() ([]*User, error) {
	users := []*User{}
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

func (udb *UserStore) Find(id string) (*User, error) {
	res, err := rethinkdb.Table(PREFIX).Get(id).Run(udb.session)
	if err != nil {
		return nil, err
	}
	if res.IsNil() {
		return nil, fmt.Errorf("User not found")
	}
	var user *User
	res.One(&user)
	res.Close()
	return user, nil
}

func (udb *UserStore) FindByField(feild, value string) (User, error) {
	res, err := rethinkdb.Table(PREFIX).Filter(rethinkdb.Row.Field(feild).Eq(value)).Run(udb.session)
	var user User
	if err != nil {
		return user, err
	}

	if res.IsNil() {
		return user, fmt.Errorf("User not found")
	}

	res.One(&user)
	res.Close()
	return user, nil
}

// func (udb *UserStore) Create(user User) error {
// 	ctx := context.Background()
// 	if _, err := udb.client.HSetNX(ctx, udb.GetRedisKey(user.ID), "id", &user).Result(); err != nil {
// 		log.Errorf("create: redis error: %w", err)
// 		return fmt.Errorf("create: redis error: %w", err)
// 	}
// 	return nil
// }

// func (udb *UserStore) Find(userID string, usePrefix bool) (*User, error) {
// 	ctx := context.Background()
// 	if usePrefix {
// 		userID = udb.GetRedisKey(userID)
// 	}
// 	result, err := udb.client.HGet(ctx, userID, "id").Result()
// 	if err != nil && err != redis.Nil {
// 		return nil, fmt.Errorf("find: redis error: %w", err)
// 	}
// 	if result == "" {
// 		return nil, fmt.Errorf("find: not found")
// 	}

// 	user := &User{}
// 	if err := user.UnmarshalBinary([]byte(result)); err != nil {
// 		return nil, fmt.Errorf("find: unmarshal error: %w", err)
// 	}

// 	return user, nil
// }

// func (udb *UserStore) GetAll() ([]*User, error) {
// 	ctx := context.Background()
// 	iter := udb.client.Scan(ctx, 0, PREFIX+"*", 0).Iterator()
// 	users := []*User{}
// 	for iter.Next(ctx) {
// 		if user, err := udb.Find(iter.Val(), false); err == nil {
// 			users = append(users, user)
// 		}
// 	}
// 	return users, nil
// }

// func (udb *UserStore) Update(ctx context.Context, user User) error {
// 	// Find token:     a.rds.HGet()
// 	// Override token: a.rds.HSet()
// 	return nil
// }

// func (udb *UserStore) Delete(ctx context.Context, tokenID string) error {
// 	// Find token:   a.rds.HGet()
// 	// Delete token: a.rds.Del()
// 	return nil
// }
