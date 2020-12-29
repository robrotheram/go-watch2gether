package room

import (
	"bytes"
	"encoding/json"

	"github.com/segmentio/ksuid"
	log "github.com/sirupsen/logrus"
)

var SERVER_USER = RoomWatcher{ID: ksuid.New().String(), Name: "Server"}

type Event struct {
	Action       string        `json:"action"`
	Host         string        `json:"host"`
	Watcher      RoomWatcher   `json:"watcher"`
	Queue        []Video       `json:"queue"`
	CurrentVideo Video         `json:"current_video"`
	Seek         float32       `json:"seek"`
	Watchers     []RoomWatcher `json:"watchers"`
	Settings     RoomSettings  `json:"settings"`
}

func (evt Event) ToBytes() []byte {
	b, _ := json.Marshal(evt)
	return b
}

func processEvent(data []byte) (Event, error) {
	var evnt Event
	in := bytes.NewReader(data)
	err := json.NewDecoder(in).Decode(&evnt)
	if err != nil {
		log.Errorf("Error Decoding Event: %v", err)
		return Event{}, err
	}
	return evnt, nil
}

func (r *Room) HandleEvent(evt Event) {
	if evt.Watcher.ID == SERVER_USER.ID {
		return
	}
	switch evt.Action {
	case EVNT_PLAYING:
		r.SetPlaying(true)
		r.SendClientEvent(evt)
	case EVNT_PAUSING:
		r.SetPlaying(false)
		r.SendClientEvent(evt)
	case EVNT_UPDATE_HOST:
		r.SetHost(evt.Host)
	case EVNT_NEXT_VIDEO:
		r.ChangeVideo()
	case EVNT_SEEK:
		r.SetSeek(evt.Seek)
	case EVNT_UPDATE_SETTINGS:
		r.SetSettings(evt.Settings)
	case EVNT_SEEK_TO_ME:
		r.SetSeek(evt.Watcher.Seek)
	case EVNT_UPDATE_QUEUE:
		r.SetQueue(evt.Queue)
	case ENVT_FINSH:
		r.HandleFinish(evt.Watcher)
	case EVNT_USER_UPDATE:
		r.SeenUser(evt.Watcher)
	case EVT_ROOM_EXIT:
		r.DeleteIfEmpty()
	}
}

func NewEvent(action string) Event {
	return Event{Action: action}
}
