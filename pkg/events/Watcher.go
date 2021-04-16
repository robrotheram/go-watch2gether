package events

import "time"

type RoomWatcher struct {
	ID       string    `json:"id"`
	Name     string    `json:"name"`
	Seek     float32   `json:"seek"`
	VideoID  string    `json:"video_id"`
	IsHost   bool      `json:"is_host"`
	LastSeen time.Time `json:"seen"`
}
