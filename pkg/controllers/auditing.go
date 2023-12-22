package controllers

import log "github.com/sirupsen/logrus"

type Auditing struct {
}

func (a *Auditing) Send(event Event) {
	log.WithFields(log.Fields{
		"user":   event.Action.User,
		"action": event.Action.Type,
	}).Info(event.Message)
}
