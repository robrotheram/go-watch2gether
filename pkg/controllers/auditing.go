package controllers

import log "github.com/sirupsen/logrus"

type Auditing struct {
}

func (a *Auditing) Send(event Event) {
	log.Infof("%s %s", event.Action.User, event.Action.ActionType)
}
