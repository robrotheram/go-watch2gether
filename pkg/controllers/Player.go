package controllers

import "w2g/pkg/media"

type Player interface {
	Play(string, int) error
	Progress() media.MediaDuration
	Pause()
	Unpause()
	Stop()
	Close()
}
