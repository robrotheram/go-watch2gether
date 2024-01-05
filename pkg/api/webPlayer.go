package api

import (
	"w2g/pkg/controllers"
	"w2g/pkg/media"
)

const WEBPLAYER = controllers.PlayerType("WEB_PLAYER")

type WebPlayer struct {
	done     chan any
	progress media.MediaDuration
	running  bool
}

func NewWebPlayer() *WebPlayer {
	return &WebPlayer{
		done: make(chan any),
		progress: media.MediaDuration{
			Progress: 0,
		},
	}
}

func (wb *WebPlayer) Type() controllers.PlayerType {
	return WEBPLAYER
}

func (wb *WebPlayer) Play(url string, start int) error {
	wb.running = true
	<-wb.done
	wb.running = false
	return nil
}

func (wb *WebPlayer) Progress() media.MediaDuration {
	return wb.progress
}

func (wb *WebPlayer) Pause()   {}
func (wb *WebPlayer) Unpause() {}

func (wb *WebPlayer) Stop() {
	if wb.running {
		wb.done <- "STOP"
	}
}
func (wb *WebPlayer) Close() {
	wb.Stop()
}

func (wb *WebPlayer) Status() bool {
	return wb.running
}

func (wb *WebPlayer) UpdateDuration(duration media.MediaDuration) {
	wb.progress = duration
}
