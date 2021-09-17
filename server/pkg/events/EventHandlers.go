package events

import (
	"watch2gether/pkg/media"
	meta "watch2gether/pkg/roomMeta"

	"github.com/prometheus/common/log"
)

type Handler = func(*Event, *meta.Meta)

var EventHandlers = map[string]Handler{
	EVNT_PLAYING:         HandlePlaying,
	EVNT_PAUSING:         HandlePause,
	EVNT_UPDATE_QUEUE:    HandleUpdateQueue,
	EVNT_USER_UPDATE:     HandleUserUpdate,
	EVNT_USER_LEAVE:      HandleUserLeave,
	EVNT_SEEK_TO_USER:    HandleForceSeek,
	ENVT_FINSH:           HandleFinish,
	EVNT_NEXT_VIDEO:      HandleNextVideo,
	EVNT_UPDATE_SETTINGS: HandleUpdateSettings,
	EVT_SUFFLE_QUEUE:     HandleSuffleQueue,
	EVNT_ADD_VIDEO:       HandleNewVideo,
}

func HandlePlaying(evt *Event, meta *meta.Meta) {
	meta.Playing = true
}

func HandleSuffleQueue(evt *Event, meta *meta.Meta) {
	meta.ShuffleQueue()
}

func HandlePause(evt *Event, meta *meta.Meta) {
	meta.Playing = false
}

func HandleUpdateQueue(evt *Event, meta *meta.Meta) {
	meta.Queue = evt.Queue
}

func HandleNextVideo(evt *Event, meta *meta.Meta) {
	meta.NextVideo()
}

func HandleForceSeek(evt *Event, meta *meta.Meta) {
	meta.SinkAllWatchers(evt.Watcher.Seek)
}

func HandleUpdateSettings(evt *Event, meta *meta.Meta) {
	meta.Settings = evt.Settings
}

func HandleUserUpdate(evt *Event, meta *meta.Meta) {
	err := meta.UpdateWatcher(evt.Watcher)
	if err != nil {
		meta.AddWatcher(evt.Watcher)
	}
}

func HandleUserLeave(evt *Event, meta *meta.Meta) {
	meta.RemoveWatcher(evt.Watcher.ID)
	if meta.Host == evt.Watcher.ID && len(meta.Watchers) > 0 {
		meta.Host = meta.Watchers[0].ID
	}
}

func HandleFinish(evt *Event, meta *meta.Meta) {
	evt.Watcher.Seek = media.SEEK_FINISHED
	meta.UpdateWatcher(evt.Watcher)

	//Do not Skip if setting is false
	if !meta.Settings.AutoSkip {
		return
	}
	// Check if a watcher has not finished
	for i := range meta.Watchers {
		u := &meta.Watchers[i]
		if !u.Seek.Done() {
			return
		}
	}
	meta.ResetWatcher()
	meta.NextVideo()
	evt.Action = EVNT_NEXT_VIDEO
	log.Debug("CHANGING VIDOE!!!!!!!!!!!!")
}

func HandleNewVideo(evt *Event, meta *meta.Meta) {
	video := evt.Video
	meta.Queue = append(meta.Queue, video)
}
