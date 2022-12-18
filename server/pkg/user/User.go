package user

import (
	"encoding/json"
	"fmt"

	"github.com/asdine/storm"
	"github.com/segmentio/ksuid"
	log "github.com/sirupsen/logrus"
)

const (
	USER_TYPE_DISCORD = "DISCORD"
	USER_TYPE_SERVER  = "SERVER"
	USER_TYPE_ANON    = "ANONYMOUS"
	USER_TYPE_BASIC   = "BASIC"
)

type User struct {
	ID         string `storm:"id" json:"id"`
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
	session *storm.DB
}

func (u *UserStore) GetKey(id string) string {
	return fmt.Sprintf("%s:%s", PREFIX, id)
}

func NewUserStore(session *storm.DB) *UserStore {
	us := &UserStore{session: session}
	us.Cleanup()
	return us
}

func (udb *UserStore) Create(user User) error {
	return udb.session.Save(user)
}

func (udb *UserStore) GetAll() ([]User, error) {
	users := []User{}
	err := udb.session.All(&users)
	return users, err
}

func (udb *UserStore) FindByType(value string) ([]User, error) {
	users, err := udb.GetAll()
	if err != nil {
		return users, err
	}
	filteredUsers := []User{}
	for _, user := range users {
		if user.Type == value {
			filteredUsers = append(filteredUsers, user)
		}
	}
	return filteredUsers, nil
}
func (udb *UserStore) FindById(value string) ([]User, error) {
	users, err := udb.GetAll()
	if err != nil {
		return users, err
	}
	filteredUsers := []User{}
	for _, user := range users {
		if user.ID == value {
			filteredUsers = append(filteredUsers, user)
		}
	}
	return filteredUsers, nil
}
func (udb *UserStore) FindByName(value string) (User, error) {
	users, err := udb.GetAll()
	if err != nil {
		return User{}, err
	}
	filteredUsers := []User{}
	for _, user := range users {
		if user.Username == value {
			filteredUsers = append(filteredUsers, user)
		}
	}
	return filteredUsers[0], nil
}
func (udb *UserStore) Delete(user User) error {
	return udb.session.DeleteStruct(user)
}

func (udb *UserStore) Cleanup() {
	users, _ := udb.FindByType(USER_TYPE_ANON)
	for _, u := range users {
		udb.Delete(u)
	}
}
