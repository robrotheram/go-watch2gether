package controllers

import log "github.com/sirupsen/logrus"

type Auditing struct {
}

func (a *Auditing) Send(event Event) {

	//Filter out update_duration is spam and not needed for the log
	if event.Action.Type == UPDATE_DURATION {
		return
	}
	log.WithFields(log.Fields{
		"channel": event.Action.Channel,
		"user":    event.Action.User,
		"action":  event.Action.Type,
	}).Info(event.Message)
}
