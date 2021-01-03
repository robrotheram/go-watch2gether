package room

import (
	"encoding/json"
	"fmt"
	"time"
	user "watch2gether/pkg/user"

	"github.com/segmentio/ksuid"
)

type RoomWatcher struct {
	ID       string    `json:"id"`
	Name     string    `json:"name"`
	Seek     float32   `json:"seek"`
	VideoID  string    `json:"video_id"`
	IsHost   bool      `json:"is_host"`
	LastSeen time.Time `json:"seen"`
}

type RoomSettings struct {
	Controls bool `json:"controls"`
	AutoSkip bool `json:"auto_skip"`
}

func NewRoomSettings() RoomSettings {
	return RoomSettings{
		Controls: true,
		AutoSkip: true,
	}

}

type Meta struct {
	ID           string        `rethinkdb:"id,omitempty" json:"id"`
	Name         string        `json:"name"`
	Owner        string        `json:"owner"`
	Host         string        `json:"host"`
	History      []Video       `json:"history"`
	CurrentVideo Video         `json:"current_video"`
	Seek         float32       `json:"seek"`
	Playing      bool          `json:"playing"`
	Queue        []Video       `json:"queue"`
	Watchers     []RoomWatcher `json:"watchers"`
	Type         string        `json:"type"`
	Settings     RoomSettings  `json:"settings"`
}

func (t *Meta) MarshalBinary() ([]byte, error) {
	return json.Marshal(t)
}

func (t *Meta) UnmarshalBinary(data []byte) error {
	if err := json.Unmarshal(data, &t); err != nil {
		return err
	}
	return nil
}

func NewMeta(name string, usr user.User) *Meta {
	if usr.Type == user.USER_TYPE_ANON {
		return &Meta{
			Name:     name,
			Watchers: []RoomWatcher{},
			Queue:    []Video{},
			History:  []Video{},
			ID:       ksuid.New().String(),
			Owner:    "",
			Host:     usr.ID,
			Type:     "Basic",
			Settings: NewRoomSettings(),
		}
	}
	return &Meta{
		Name:     name,
		Watchers: []RoomWatcher{},
		Queue:    []Video{},
		History:  []Video{},
		ID:       ksuid.New().String(),
		Owner:    usr.ID,
		Host:     usr.ID,
		Type:     "Basic",
		Settings: NewRoomSettings(),
	}

}
func (meta *Meta) GetLastVideo() Video {
	if len(meta.History) == 0 {
		return Video{}
	}
	return meta.History[len(meta.History)-1]
}

func (meta *Meta) UpdateHistory(v Video) {
	if v.ID == "" {
		return
	}
	meta.History = append(meta.History, v)
}

func (meta *Meta) FindWatcher(id string) (RoomWatcher, error) {
	for _, w := range meta.Watchers {
		if w.ID == id {
			return w, nil
		}
	}
	return RoomWatcher{}, fmt.Errorf("Watcher not found")
}
func (meta *Meta) RemoveWatcher(watcherid string) {
	for i, v := range meta.Watchers {
		if v.ID == watcherid {
			meta.Watchers = append(meta.Watchers[:i], meta.Watchers[i+1:]...)
			break
		}
	}
}

func (meta *Meta) AddWatcher(rw RoomWatcher) {
	if _, err := meta.FindWatcher(rw.ID); err == nil {
		return
	}
	rw.LastSeen = time.Now()
	meta.Watchers = append(meta.Watchers, rw)
}

func (meta *Meta) UpdateWatcher(rw RoomWatcher) error {
	for i := range meta.Watchers {
		watcher := &meta.Watchers[i]
		if rw.ID == watcher.ID {
			if watcher.IsHost {
				meta.Seek = rw.Seek
			}
			watcher.LastSeen = time.Now()
			watcher.Seek = rw.Seek
			watcher.VideoID = rw.VideoID

			return nil
		}
	}
	return fmt.Errorf("Watcher not found")
}

func NewWatcher(usr user.User) RoomWatcher {
	return RoomWatcher{ID: usr.ID, Name: usr.Name}
}
