package pkg

import "time"

type User struct {
	Name     string  `json:"name"`
	Seek     float32 `json:"seek"`
	IsHost   bool    `json:"false"`
	LastSeen time.Time
}

func NewUser(name string) User {
	return User{
		Name:     name,
		IsHost:   false,
		Seek:     0,
		LastSeen: time.Now(),
	}
}
