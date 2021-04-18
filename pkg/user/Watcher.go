package user

import (
	"time"
)

type Watcher struct {
	User
	Seek     float64   `json:"seek"`
	VideoID  string    `json:"video_id"`
	IsHost   bool      `json:"is_host"`
	LastSeen time.Time `json:"seen"`
}

func NewWatcher(usr User) Watcher {
	return Watcher{User: usr}
}

var SERVER_USER = Watcher{
	User: User{ID: "cec98f5d-ea6a-414c-ad1c-331cad9f01af", Username: "Server", Type: USER_TYPE_SERVER},
}

var DISCORD_BOT = Watcher{
	User: User{ID: "668bef28-6f2f-40c7-a0fa-65f9d06cd086", Username: "Bot", Type: USER_TYPE_DISCORD},
}
