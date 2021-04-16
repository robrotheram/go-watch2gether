package user

import (
	"encoding/json"
	"fmt"

	"github.com/segmentio/ksuid"
	log "github.com/sirupsen/logrus"
	rethinkdb "gopkg.in/rethinkdb/rethinkdb-go.v6"
)

const (
	USER_TYPE_DISCORD = "DISCORD"
	USER_TYPE_ANON    = "ANONYMOUS"
	USER_TYPE_BASIC   = "BASIC"
)

type User struct {
	ID         string `rethinkdb:"id,omitempty" json:"id"`
	Username   string `json:"username"`
	Type       string `json:"type"`
	Avatar     string `json:"avatar"`
	AvatarIcon string `json:"avatar_icon"`
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

func NewUser(name string, _type string) User {
	id := ksuid.New().String()
	log.Infof("Creating User with ID %s, Name: %s", id, name)
	return User{
		ID:       id,
		Username: name,
		Type:     _type,
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
	us := &UserStore{session: session}
	us.Cleanup()
	return us
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

func (udb *UserStore) FindAllByField(feild, value string) ([]User, error) {
	res, err := rethinkdb.Table(PREFIX).Filter(rethinkdb.Row.Field(feild).Eq(value)).Run(udb.session)
	users := []User{}
	if err != nil {
		return users, err
	}

	if res.IsNil() {
		return users, fmt.Errorf("User not found")
	}

	res.All(&users)
	res.Close()

	return users, nil
}

func (udb *UserStore) Delete(id string) error {
	_, err := rethinkdb.Table(PREFIX).Get(id).Delete().RunWrite(udb.session)
	return err
}

func (udb *UserStore) Cleanup() {
	users, _ := udb.FindAllByField("Type", USER_TYPE_ANON)
	for _, u := range users {
		udb.Delete(u.ID)
	}
}
