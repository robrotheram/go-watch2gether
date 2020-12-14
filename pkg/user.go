package pkg

import "time"

type User struct {
	Name         string  `json:"name"`
	Seek         float32 `json:"seek"`
	IsHost       bool    `json:"is_host"`
	CurrentVideo Video   `json:"current_video"`
	LastSeen     time.Time
}

func NewUser(name string) User {
	return User{
		Name:     name,
		IsHost:   false,
		Seek:     0,
		LastSeen: time.Now(),
	}
}
