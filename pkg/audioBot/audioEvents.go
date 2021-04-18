package audioBot

import (
	"watch2gether/pkg/events"
	"watch2gether/pkg/user"
)

func CreateBotJoinEvent() events.Event {
	return events.Event{
		Action:  events.EVNT_USER_UPDATE,
		Watcher: user.DISCORD_BOT,
	}
}

func CreateBotUpdateEvent(seek float64) events.Event {
	evt := events.Event{
		Action:  events.EVNT_USER_UPDATE,
		Watcher: user.DISCORD_BOT,
	}
	evt.Watcher.Seek = seek
	return evt
}

func CreateBotFinishEvent() events.Event {
	evt := events.Event{
		Action:  events.ENVT_FINSH,
		Watcher: user.DISCORD_BOT,
	}
	evt.Watcher.Seek = float64(1)
	return evt
}
