package audioBot

import (
	"fmt"
	"sync"
	"time"
	"watch2gether/pkg/events"
	"watch2gether/pkg/media"
	"watch2gether/pkg/utils"

	"github.com/bwmarrin/discordgo"
	log "github.com/sirupsen/logrus"
)

type AudioBot struct {
	session              *discordgo.Session
	audio                *Audio
	voiceChannelID       string
	notficationChannelID string
	VoiceConnection      *discordgo.VoiceConnection
	RoomChannel          chan []byte
	Running              bool
	updateTime           time.Time
	ticker               *time.Ticker
	sync.Mutex
}

func NewAudioBot(voiceCh string, notificationCh string, vc *discordgo.VoiceConnection, s *discordgo.Session) *AudioBot {
	ab := AudioBot{
		session:              s,
		voiceChannelID:       voiceCh,
		notficationChannelID: notificationCh,
		VoiceConnection:      vc,
		Running:              false,
	}
	return &ab
}

func (ab *AudioBot) SendToRoom(evt events.Event) {
	if ab.Running {
		ab.RoomChannel <- evt.ToBytes()
	}
}

func (ab *AudioBot) RegisterToRoom(rc chan []byte) {
	ab.RoomChannel = rc
	ab.audio = NewAudio(ab, ab.VoiceConnection)
	ab.Running = true

}

func (ab *AudioBot) Send(evt events.RoomState) {
	ab.Lock()
	defer ab.Unlock()
	ab.handleEvent(evt)
}

func (ab *AudioBot) sendToChannel(msg string) {
	ab.session.ChannelMessageSend(ab.notficationChannelID, msg)
}

func (ab *AudioBot) handleEvent(evt events.RoomState) {
	if ab.audio == nil {
		ab.audio = NewAudio(ab, ab.VoiceConnection)
	}
	switch evt.Action {
	case events.EVNT_UPDATE_QUEUE:
		if utils.Configuration.DiscordNotify {
			ab.sendToChannel(fmt.Sprintf("Queue Updated by: %s", evt.Watcher.Username))
		}
	case events.EVNT_NEXT_VIDEO:
		if evt.Playing {
			ab.PlayAudio(evt.CurrentVideo, 0)
		} else {
			ab.audio.Stop()
		}
	case events.EVNT_PLAYING:
		if !ab.audio.Playing {
			ab.PlayAudio(evt.Meta.CurrentVideo, 0)
		} else {
			ab.audio.Unpause()
		}
		if utils.Configuration.DiscordNotify {
			ab.sendToChannel(fmt.Sprintf("User: %s Started the video", evt.Watcher.Username))
		}
	case events.EVNT_PAUSING:
		ab.audio.Paused()
		if utils.Configuration.DiscordNotify {
			ab.sendToChannel(fmt.Sprintf("User: %s Paused the video", evt.Watcher.Username))
		}
	case events.EVNT_SEEK_TO_USER:
		ab.audio.Stop()
		ab.PlayAudio(evt.CurrentVideo, int(evt.GetHostSeek().ProgressSec))
		if utils.Configuration.DiscordNotify {
			ab.sendToChannel(fmt.Sprintf("User: %s Seeked Video to %f", evt.Watcher.Username, evt.GetHostSeek().ProgressPct*100))
		}
	case events.EVT_ROOM_EXIT:
		ab.sendToChannel(("room has closed down"))
	}
}

func (ab *AudioBot) PlayAudio(video media.Video, starttime int) {
	ab.updateTime = time.Now()
	switch video.GetType() {
	case media.VIDEO_TYPE_YT:
		ab.PlayYoutube(video.Url, starttime)
	case media.VIDEO_TYPE_MP4:
		ab.PlayAudioFile(video.Url, starttime)
	default:
		if ab.audio != nil {
			ab.audio.Stop()
			ab.audio = nil
		}
		log.Debugf("Video Type could not be found %v", video)
	}
}

func (ab *AudioBot) PlayYoutube(videoURL string, starttime int) {

	downloadURL, err := media.GetYoutubeURL(videoURL)
	if err != nil {
		log.Warnf("unable to get youtube url: %v", err)
		return
	}
	ab.PlayAudioFile(downloadURL, starttime)
}
func (ab *AudioBot) PlayAudioFile(url string, starttime int) {
	if ab.audio == nil {
		log.Info("Bot not connected to Room")
	}
	if ab.audio.Playing {
		ab.audio.Stop()
		log.Debug("Stopping Audio")
	}
	err := ab.audio.Play(url, starttime)
	if err != nil {
		log.Debugf("Error Encoding URL : %v \n", err)
		return
	}
	err = ab.audio.Start()
	if err != nil {
		fmt.Printf("Error Starting Audio : %v \n", err)
		return
	}
}

func (ab *AudioBot) Disconnect() error {
	ab.SendToRoom(CreateBotLeaveEvent())
	ab.Running = false
	if ab.audio != nil {
		ab.audio.Stop()
	}
	return ab.VoiceConnection.Disconnect()
}
