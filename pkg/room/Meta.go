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

type Meta struct {
	ID           string        `rethinkdb:"id,omitempty" json:"id"`
	Name         string        `json:"name"`
	Owner        string        `json:"owner"`
	Host         string        `json:"host"`
	History      []Video       `json:"history"`
	CurrentVideo Video         `json:"current_video"`
	Seek         float32       `json:"seek"`
	Controls     bool          `json:"controls"`
	Playing      bool          `json:"playing"`
	Queue        []Video       `json:"queue"`
	Watchers     []RoomWatcher `json:"watchers"`
	Type         string        `json:"type"`
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

func NewMeta(name string, owner string) *Meta {
	return &Meta{
		Name:     name,
		Watchers: []RoomWatcher{},
		Queue:    []Video{},
		History:  []Video{},
		Controls: true,
		ID:       ksuid.New().String(),
		Owner:    owner,
		Host:     owner,
		Type:     "Basic",
	}

}
func (meta *Meta) GetLastVideo() Video {
	if len(meta.History) == 0 {
		return Video{}
	}
	return meta.History[len(meta.History)-1]
}

func (meta *Meta) UpdateHistory(v Video) {
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

func (m *Meta) Update(meta Meta) {

	if m.Name != meta.Name && meta.Name != "" {
		m.Name = meta.Name
	}
	if m.Host != meta.Host && meta.Host != "" {
		m.Host = meta.Host
	}
	if m.CurrentVideo != meta.CurrentVideo && meta.CurrentVideo.Url != "" {
		m.CurrentVideo = meta.CurrentVideo
	}
	if m.Seek != meta.Seek {
		m.Seek = meta.Seek
	}
	if m.Controls != meta.Controls {
		m.Controls = meta.Controls
	}
	if m.Playing != meta.Playing {
		m.Playing = meta.Playing
	}
	if meta.Queue != nil {
		m.Queue = meta.Queue
	}
}
