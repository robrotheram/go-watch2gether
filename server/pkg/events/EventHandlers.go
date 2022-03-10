package events

import (
	"watch2gether/pkg/media"
	meta "watch2gether/pkg/roomMeta"

	log "github.com/sirupsen/logrus"
)

type Handler = func(*Event, *meta.Meta)

var EventHandlers = map[string]Handler{
	EVENT_PLAYING:         HandlePlaying,
	EVENT_PAUSING:         HandlePause,
	EVENT_UPDATE_QUEUE:    HandleUpdateQueue,
	EVENT_USER_UPDATE:     HandleUserUpdate,
	EVENT_USER_LEAVE:      HandleUserLeave,
	EVENT_SEEK_TO_USER:    HandleForceSeek,
	EVENT_FINISH:          HandleFinish,
	EVENT_NEXT_VIDEO:      HandleNextVideo,
	EVENT_UPDATE_SETTINGS: HandleUpdateSettings,
	EVENT_SHUFFLE_QUEUE:   HandleShuffleQueue,
	EVENT_ADD_VIDEO:       HandleNewVideo,
	EVENT_BOT_LEAVE:       EmptyHandler,
	EVENT_ROOM_EXIT:       EmptyHandler,
}

func EmptyHandler(evt *Event, meta *meta.Meta) {}

func HandlePlaying(evt *Event, meta *meta.Meta) {
	meta.Playing = true
}

func HandleShuffleQueue(evt *Event, meta *meta.Meta) {
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
	evt.Action = EVENT_NEXT_VIDEO
	log.Debug("CHANGING VIDEO!!!!!!!!!!!!")
}

func HandleNewVideo(evt *Event, meta *meta.Meta) {
	video := evt.Video
	meta.Queue = append(meta.Queue, video)
}
