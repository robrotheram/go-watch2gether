package room

import (
	"encoding/json"
	"fmt"
	"time"
	events "watch2gether/pkg/events"
	"watch2gether/pkg/media"
	user "watch2gether/pkg/user"

	"github.com/segmentio/ksuid"
)

type Meta struct {
	ID           string               `rethinkdb:"id,omitempty" json:"id"`
	Name         string               `json:"name"`
	Owner        string               `json:"owner"`
	Host         string               `json:"host"`
	Icon         string               `json:"icon"`
	History      []media.Video        `json:"history"`
	CurrentVideo media.Video          `json:"current_video"`
	Seek         float32              `json:"seek"`
	Playing      bool                 `json:"playing"`
	Queue        []media.Video        `json:"queue"`
	Watchers     []events.RoomWatcher `json:"watchers"`
	Type         string               `json:"type"`
	Settings     events.RoomSettings  `json:"settings"`
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
			Watchers: []events.RoomWatcher{},
			Queue:    []media.Video{},
			History:  []media.Video{},
			ID:       ksuid.New().String(),
			Owner:    "",
			Host:     usr.ID,
			Type:     "Basic",
			Settings: events.NewRoomSettings(),
		}
	}
	return &Meta{
		Name:     name,
		Watchers: []events.RoomWatcher{},
		Queue:    []media.Video{},
		History:  []media.Video{},
		ID:       ksuid.New().String(),
		Owner:    usr.ID,
		Host:     usr.ID,
		Type:     "Basic",
		Settings: events.NewRoomSettings(),
	}

}
func (meta *Meta) GetLastVideo() media.Video {
	if len(meta.History) == 0 {
		return media.Video{}
	}
	return meta.History[len(meta.History)-1]
}

func (meta *Meta) UpdateHistory(v media.Video) {
	if v.ID == "" {
		return
	}
	meta.History = append(meta.History, v)
}

func (meta *Meta) FindWatcher(id string) (events.RoomWatcher, error) {
	for _, w := range meta.Watchers {
		if w.ID == id {
			return w, nil
		}
	}
	return events.RoomWatcher{}, fmt.Errorf("Watcher not found")
}
func (meta *Meta) RemoveWatcher(watcherid string) {
	for i, v := range meta.Watchers {
		if v.ID == watcherid {
			meta.Watchers = append(meta.Watchers[:i], meta.Watchers[i+1:]...)
			break
		}
	}
}

func (meta *Meta) AddWatcher(rw events.RoomWatcher) {
	if _, err := meta.FindWatcher(rw.ID); err == nil {
		return
	}
	rw.LastSeen = time.Now()
	meta.Watchers = append(meta.Watchers, rw)
}

func (meta *Meta) UpdateWatcher(rw events.RoomWatcher) error {
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

func NewWatcher(usr user.User) events.RoomWatcher {
	return events.RoomWatcher{ID: usr.ID, Name: usr.Username}
}
