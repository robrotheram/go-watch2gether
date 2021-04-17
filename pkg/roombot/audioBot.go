package roombot

import (
	"fmt"
	"sync"
	"watch2gether/pkg/events"
	"watch2gether/pkg/media"

	"github.com/bwmarrin/discordgo"
	"github.com/jonas747/dca"
	"github.com/kkdai/youtube/v2"
)

type AudioBot struct {
	session              *discordgo.Session
	audio                *Audio
	voiceChannelID       string
	notficationChannelID string
	VoiceConnection      *discordgo.VoiceConnection
	EventChannel         chan events.Event
	roomChannel          chan events.Event
	stream               *dca.StreamingSession
	Running              bool
	wg                   sync.WaitGroup
	sync.Mutex
}

func NewAudioBot(voiceCh string, notificationCh string, vc *discordgo.VoiceConnection, s *discordgo.Session) *AudioBot {
	ab := AudioBot{
		session:              s,
		voiceChannelID:       voiceCh,
		notficationChannelID: notificationCh,
		VoiceConnection:      vc,
		EventChannel:         make(chan events.Event),
		Running:              false,
	}
	return &ab
}

func (ab *AudioBot) Start() error {
	ab.Lock()
	start := !ab.Running
	ab.Running = true
	ab.Unlock()
	if start {
		ab.wg.Add(1)
		go func() {
			ab.Watcher()
			ab.Lock()
			ab.Running = false
			defer ab.wg.Done()
			ab.Unlock()
		}()
		return nil
	} else {
		return fmt.Errorf("Bot already started")
	}

}

func (ab *AudioBot) Stop() {
	if ab.Running {
		close(ab.EventChannel)
		ab.wg.Wait()
	}
}

func (ab *AudioBot) Send(evt events.Event) {
	ab.EventChannel <- evt
}

func (ab *AudioBot) Watcher() {
	for evt := range ab.EventChannel {
		ab.handleEvent(evt)
	}
}

func (ab *AudioBot) sendToChannel(msg string) {
	ab.session.ChannelMessageSend(ab.notficationChannelID, msg)
}

func (ab *AudioBot) handleEvent(evt events.Event) {
	switch evt.Action {
	case events.EVNT_UPDATE_QUEUE:
		ab.sendToChannel(fmt.Sprintf("Queue Updated by: %s", evt.Watcher.Username))
	case events.EVT_VIDEO_CHANGE:
		ab.PlayAudio(evt.CurrentVideo)
	case events.EVT_ROOM_EXIT:
		ab.sendToChannel(fmt.Sprintf("Room has closed down"))
	case events.EVNT_PLAYING:
		if ab.audio == nil {
			ab.PlayAudio(evt.CurrentVideo)
		} else {
			ab.audio.Unpause()
		}
		ab.sendToChannel(fmt.Sprintf("User: %s Started the video", evt.Watcher.Username))
	case events.EVNT_PAUSING:
		ab.audio.Paused()
		ab.sendToChannel(fmt.Sprintf("User: %s Paused the video", evt.Watcher.Username))
	}
}

func (ab *AudioBot) PlayAudio(video media.Video) {
	switch video.GetType() {
	case media.VIDEO_TYPE_YT:
		ab.PlayYoutube(video.Url)
		break
	case media.VIDEO_TYPE_MP4:
		ab.PlayAudioFile(video.Url)
		break
	default:
		ab.audio.Stop()
	}
}

func (ab *AudioBot) PlayYoutube(videoURL string) {
	client := youtube.Client{}
	video, _ := client.GetVideo(videoURL)
	downloadURL, _ := client.GetStreamURL(video, &video.Formats[0])
	ab.PlayAudioFile(downloadURL)
}

func (ab *AudioBot) PlayAudioFile(url string) {
	if ab.audio != nil {
		ab.audio.Stop()
	}
	audio, err := NewAudio(url, ab.VoiceConnection)
	if err != nil {
		fmt.Println("ERROR!!!!")
	}
	ab.audio = audio
	ab.audio.Play()
}
