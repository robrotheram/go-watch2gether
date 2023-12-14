package model

import "watch2gether/pkg/media"

type Action string

type Event struct {
	Action  Action
	Message string
	media.Media
}

const (
	PLAY        = Action("PLAY")
	PAUSE       = Action("PAUSE")
	SKIP        = Action("SKIP")
	UPDATEQUEUE = Action("UPDATE_QUEUE")
	SHUFFLE     = Action("SHUFFLE")
	STOP        = Action("STOP")
	UPADATE     = Action("UPDATE")
	FINISHED    = Action("FINISHED")
	MESSAGE     = Action("MESSAGE")
)

func NewEvent(m media.Media) Event {
	return Event{
		Media: m,
	}
}

func (event Event) WithAction(action Action) Event {
	event.Action = action
	return event
}

func (event *Event) WithMsg(msg string) *Event {
	event.Message = msg
	return event
}
