package user

import (
	"time"

	"github.com/segmentio/ksuid"
)

type Watcher struct {
	User
	Seek     float32   `json:"seek"`
	VideoID  string    `json:"video_id"`
	IsHost   bool      `json:"is_host"`
	LastSeen time.Time `json:"seen"`
}

func NewWatcher(usr User) Watcher {
	return Watcher{User: usr}
}

var SERVER_USER = Watcher{
	User: User{ID: ksuid.New().String(), Username: "Server"},
}

var DISCORD_BOT = Watcher{
	User: User{ID: ksuid.New().String(), Username: "Bot"},
}
