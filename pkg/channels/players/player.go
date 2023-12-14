package players

import (
	"time"
	"watch2gether/pkg/channels/model"
)

type PlayerType string

type Player interface {
	Start(chan model.Event)
	Play(model.Event) error
	Pause() error
	Stop() error
	Quit() error
	Duration() time.Duration
	Notify(model.Event)
}
