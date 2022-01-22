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
	Audio                *Audio
	GuildID              string
	voiceChannelID       string
	notficationChannelID string
	VoiceConnection      *discordgo.VoiceConnection
	RoomChannel          chan []byte
	Running              bool
	updateTime           time.Time
	ticker               *time.Ticker
	done                 chan bool
	sync.Mutex
}

func NewAudioBot(voiceCh string, notificationCh string, guildID string, vc *discordgo.VoiceConnection, s *discordgo.Session) *AudioBot {
	ab := AudioBot{
		session:              s,
		voiceChannelID:       voiceCh,
		notficationChannelID: notificationCh,
		VoiceConnection:      vc,
		GuildID:              guildID,
		Running:              false,
		ticker:               time.NewTicker(time.Second * 10),
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
	ab.Audio = NewAudio(ab, ab.VoiceConnection)
	ab.Running = true
	ab.Start()

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
	if ab.Audio == nil {
		ab.Audio = NewAudio(ab, ab.VoiceConnection)
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
			ab.Audio.Stop()
		}
	case events.EVNT_PLAYING:
		if !ab.Audio.Playing {
			ab.PlayAudio(evt.Meta.CurrentVideo, 0)
		} else {
			ab.Audio.Unpause()
		}
		if utils.Configuration.DiscordNotify {
			ab.sendToChannel(fmt.Sprintf("User: %s Started the video", evt.Watcher.Username))
		}
	case events.EVNT_PAUSING:
		ab.Audio.Paused()
		if utils.Configuration.DiscordNotify {
			ab.sendToChannel(fmt.Sprintf("User: %s Paused the video", evt.Watcher.Username))
		}
	case events.EVNT_SEEK_TO_USER:
		ab.Audio.Stop()
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
	case media.VIDEO_TYPE_MP3:
		ab.PlayAudioFile(video.Url, starttime)
	default:
		if ab.Audio != nil {
			ab.Audio.Stop()
			ab.Audio = nil
		}
		log.Debugf("Video Type could not be found %v", video.Type)
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
	if ab.Audio == nil {
		log.Info("Bot not connected to Room")
	}
	if ab.Audio.Playing {
		ab.Audio.Stop()
		log.Debug("Stopping Audio")
	}
	err := ab.Audio.Play(url, starttime)
	if err != nil {
		log.Debugf("Error Encoding URL : %v \n", err)
		return
	}
	err = ab.Audio.Start()
	if err != nil {
		fmt.Printf("Error Starting Audio : %v \n", err)
		return
	}
}

func (ab *AudioBot) Disconnect() error {
	log.Info("Bot diconnecting")
	ab.SendToRoom(CreateBotLeaveEvent())
	if ab.Audio != nil {
		ab.Audio.Stop()
	}
	err := ab.VoiceConnection.Disconnect()
	ab.SendToRoom(events.Event{Action: events.EVNT_BOT_LEAVE})
	ab.Running = false
	log.Error(err)
	return err
}

/**
* Function that checks if we are the only one left in a channel and disconnect if so
 */
func (ab *AudioBot) Start() {
	go func() {
		for {
			select {
			case <-ab.done:
				fmt.Println("DONE")
				return
			case <-ab.ticker.C:
				ab.LeaveCheck()
			}
		}
	}()
}
func (ab *AudioBot) Stop() {
	ab.ticker.Stop()
}

func (ab *AudioBot) LeaveCheck() {
	for _, g := range ab.session.State.Guilds {
		if g.ID == ab.GuildID {
			if len(g.VoiceStates) <= 1 {
				ab.Disconnect()
			}
		}
	}
}
