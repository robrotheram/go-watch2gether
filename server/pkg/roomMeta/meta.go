package meta

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"time"
	"watch2gether/pkg/media"
	user "watch2gether/pkg/user"

	"github.com/segmentio/ksuid"
)

type Meta struct {
	ID           string         `storm:"id" json:"id"`
	Name         string         `json:"name"`
	Owner        string         `json:"owner"`
	Host         string         `json:"host"`
	Icon         string         `json:"icon"`
	CurrentVideo media.Media    `json:"current_video"`
	Playing      bool           `json:"playing"`
	Queue        []media.Media  `json:"queue"`
	Watchers     []user.Watcher `json:"watchers"`
	Type         string         `json:"type"`
	Settings     RoomSettings   `json:"settings"`
}

func (t *Meta) MarshalBinary() []byte {
	b, _ := json.Marshal(t)
	return b
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
			Watchers: []user.Watcher{},
			Queue:    []media.Media{},
			ID:       ksuid.New().String(),
			Owner:    "",
			Host:     usr.ID,
			Type:     "Basic",
			Settings: NewRoomSettings(),
		}
	}
	return &Meta{
		Name:     name,
		Watchers: []user.Watcher{},
		Queue:    []media.Media{},
		ID:       ksuid.New().String(),
		Owner:    usr.ID,
		Host:     usr.ID,
		Type:     "Basic",
		Settings: NewRoomSettings(),
	}
}
func (meta *Meta) GetHostSeek() media.Seek {
	w, err := meta.FindWatcher(meta.Host)
	if err != nil {
		return media.SEEK_INIT
	}
	return w.Seek
}

func (meta *Meta) NextVideo() {
	if len(meta.Queue) > 0 {
		meta.CurrentVideo = meta.Queue[0]
		meta.Queue = meta.Queue[1:]
	} else {
		meta.CurrentVideo = media.Media{}
	}
}

func (meta *Meta) ShuffleQueue() {
	rand.Shuffle(len(meta.Queue), func(i, j int) {
		meta.Queue[i], meta.Queue[j] = meta.Queue[j], meta.Queue[i]
	})
}

func (meta *Meta) FindWatcher(id string) (user.Watcher, error) {
	for _, w := range meta.Watchers {
		if w.ID == id {
			return w, nil
		}
	}
	return user.Watcher{}, fmt.Errorf("watcher not found")
}

func (meta *Meta) RemoveWatcher(watcherid string) {
	for i, v := range meta.Watchers {
		if v.ID == watcherid {
			meta.Watchers = append(meta.Watchers[:i], meta.Watchers[i+1:]...)
			break
		}
	}

	if watcherid == meta.Host {
		if len(meta.Watchers) > 0 {
			meta.Host = meta.Watchers[0].ID
		} else {
			meta.Host = ""
		}
	}
}

func (meta *Meta) ResetWatcher() {
	meta.SinkAllWatchers(media.SEEK_INIT)
}
func (meta *Meta) SinkAllWatchers(seek media.Seek) {
	for i := range meta.Watchers {
		meta.Watchers[i].Seek = seek
	}
}

func (meta *Meta) AddWatcher(rw user.Watcher) {
	if _, err := meta.FindWatcher(rw.ID); err == nil {
		return
	}
	rw.LastSeen = time.Now()
	meta.Watchers = append(meta.Watchers, rw)
}

func (meta *Meta) UpdateWatcher(rw user.Watcher) error {
	for i := range meta.Watchers {
		watcher := &meta.Watchers[i]
		if rw.ID == watcher.ID {
			watcher.LastSeen = time.Now()
			watcher.Seek = rw.Seek
			watcher.VideoID = rw.VideoID
			return nil
		}
	}
	return fmt.Errorf("watcher not found")
}

func insert(array []media.Media, value media.Media, index int) []media.Media {
	return append(array[:index], append([]media.Media{value}, array[index:]...)...)
}

func remove(array []media.Media, index int) []media.Media {
	return append(array[:index], array[index+1:]...)
}

func move(array []media.Media, srcIndex int, dstIndex int) []media.Media {
	value := array[srcIndex]
	return insert(remove(array, srcIndex), value, dstIndex)
}

func (meta *Meta) ReorderQueue(pos1 int, pos2 int) {
	meta.Queue = move(meta.Queue, pos1, pos2)
}
func (meta *Meta) RemoveFromQueue(pos1 int) {
	meta.Queue = remove(meta.Queue, pos1)
}
