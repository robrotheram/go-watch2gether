package events

import (
	"bytes"
	"encoding/json"
	"watch2gether/pkg/media"

	"github.com/segmentio/ksuid"
	log "github.com/sirupsen/logrus"
)

var SERVER_USER = RoomWatcher{ID: ksuid.New().String(), Name: "Server"}

type Event struct {
	Action       string        `json:"action"`
	Host         string        `json:"host"`
	Watcher      RoomWatcher   `json:"watcher"`
	Queue        []media.Video `json:"queue"`
	CurrentVideo media.Video   `json:"current_video"`
	Seek         float32       `json:"seek"`
	Watchers     []RoomWatcher `json:"watchers"`
	Settings     RoomSettings  `json:"settings"`
}

func (evt Event) ToBytes() []byte {
	b, _ := json.Marshal(evt)
	return b
}

func ProcessEvent(data []byte) (Event, error) {
	var evnt Event
	in := bytes.NewReader(data)
	err := json.NewDecoder(in).Decode(&evnt)
	if err != nil {
		log.Errorf("Error Decoding Event: %v", err)
		return Event{}, err
	}
	return evnt, nil
}

func NewEvent(action string) Event {
	return Event{Action: action}
}
