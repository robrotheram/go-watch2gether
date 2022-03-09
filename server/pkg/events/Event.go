package events

import (
	"bytes"
	"encoding/json"
	"fmt"
	"watch2gether/pkg/media"
	meta "watch2gether/pkg/roomMeta"
	"watch2gether/pkg/user"

	log "github.com/sirupsen/logrus"
)

type Event struct {
	Action   string            `json:"action"`
	Playing  bool              `json:"playing"`
	Watcher  user.Watcher      `json:"watcher"`
	Queue    []media.Media     `json:"queue"`
	Video    media.Media       `json:"video"`
	Seek     media.Seek        `json:"seek"`
	Settings meta.RoomSettings `json:"settings"`
}

type RoomState struct {
	meta.Meta
	Action  string       `json:"action"`
	Watcher user.Watcher `json:"watcher"`
}

func (evt RoomState) ToBytes() []byte {
	b, _ := json.Marshal(evt)
	return b
}

func (evt Event) ToBytes() []byte {
	b, _ := json.Marshal(evt)
	return b
}
func (e *Event) Handle(meta *meta.Meta) (RoomState, error) {

	handler, found := EventHandlers[e.Action]
	log.Debugf("%s : handleing Event: %s", e.Watcher.Username, e.Action)
	if !found {
		log.Warnf("handler %s not found", e.Action)
		return RoomState{}, fmt.Errorf("handler %s not found", e.Action)
	}
	handler(e, meta)
	return RoomState{Meta: *meta, Watcher: e.Watcher, Action: e.Action}, nil
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
