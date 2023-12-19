package api

import "w2g/pkg/media"

type WebPlayer struct {
	done     chan any
	progress media.MediaDuration
}

func NewWebPlayer() *WebPlayer {
	return &WebPlayer{
		done: make(chan any),
		progress: media.MediaDuration{
			Progress: 0,
		},
	}
}

func (wb *WebPlayer) Play(string, int) error {
	<-wb.done
	return nil
}
func (wb *WebPlayer) Progress() media.MediaDuration {
	return wb.progress
}
func (wb *WebPlayer) Pause()   {}
func (wb *WebPlayer) Unpause() {}
func (wb *WebPlayer) Stop() {
	wb.done <- "STOP"
}
func (wb *WebPlayer) Close() {}

func (wb *WebPlayer) UpdateDuration(duration media.MediaDuration) {
	wb.progress = duration
}
