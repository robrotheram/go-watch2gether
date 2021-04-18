package audioBot

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
	RoomChannel          chan []byte
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
		Running:              false,
	}
	return &ab
}

func SendToChannel(evt events.Event, roomChannel chan []byte) {
	roomChannel <- evt.ToBytes()
}

// func (ab *AudioBot) Start() error {
// 	ab.Lock()
// 	start := !ab.Running
// 	ab.Running = true
// 	ab.Unlock()
// 	if start {
// 		ab.wg.Add(1)
// 		go func() {
// 			SendToChannel(CreateBotJoinEvent(), ab.RoomChannel)
// 			ab.Watcher()
// 			ab.Lock()
// 			ab.Running = false
// 			defer ab.wg.Done()
// 			ab.Unlock()
// 		}()
// 		return nil
// 	} else {
// 		return fmt.Errorf("Bot already started")
// 	}

// }

// func (ab *AudioBot) Stop() {
// 	if ab.Running {
// 		ab.wg.Wait()
// 	}
// }

func (ab *AudioBot) Send(evt events.Event) {
	go func() { ab.handleEvent(evt) }()
}

func (ab *AudioBot) sendToChannel(msg string) {
	ab.session.ChannelMessageSend(ab.notficationChannelID, msg)
}

func (ab *AudioBot) handleEvent(evt events.Event) {
	fmt.Println(evt.Action)
	switch evt.Action {
	case events.EVNT_UPDATE_QUEUE:
		ab.sendToChannel(fmt.Sprintf("Queue Updated by: %s", evt.Watcher.Username))

	case events.EVT_VIDEO_CHANGE:
		ab.PlayAudio(evt.CurrentVideo, int(evt.Seek.ProgressSec))
		if !evt.Playing {
			if ab.audio != nil {
				ab.audio.Paused()
			}
		}
	case events.EVNT_PLAYING:
		if ab.audio == nil {
			ab.PlayAudio(evt.CurrentVideo, int(evt.Seek.ProgressSec))
		} else {
			ab.audio.Unpause()
		}
		ab.sendToChannel(fmt.Sprintf("User: %s Started the video", evt.Watcher.Username))
	case events.EVNT_PAUSING:
		if ab.audio != nil {
			ab.audio.Paused()
		}
		ab.sendToChannel(fmt.Sprintf("User: %s Paused the video", evt.Watcher.Username))

	case events.EVNT_SEEK:
		ab.audio.Stop()
		ab.audio = nil
		ab.PlayAudio(evt.CurrentVideo, int(evt.Seek.ProgressSec))
		ab.sendToChannel(fmt.Sprintf("User: %s Seeked Video to %f%", evt.Watcher.Username, evt.Seek.ProgressPct*100))

	case events.EVT_ROOM_EXIT:
		ab.sendToChannel(fmt.Sprintf("Room has closed down"))
	}
}

func (ab *AudioBot) PlayAudio(video media.Video, starttime int) {
	fmt.Println(video)
	switch video.GetType() {
	case media.VIDEO_TYPE_YT:
		ab.PlayYoutube(video.Url, starttime)
		break
	case media.VIDEO_TYPE_MP4:
		ab.PlayAudioFile(video.Url, starttime)
		break
	default:
		if ab.audio != nil {
			ab.audio.Stop()
			ab.audio = nil
		}
		fmt.Printf("Video Type could not be found %v", video)
	}
}

func (ab *AudioBot) PlayYoutube(videoURL string, starttime int) {
	client := youtube.Client{}
	video, _ := client.GetVideo(videoURL)
	downloadURL, _ := client.GetStreamURL(video, &video.Formats[0])
	ab.PlayAudioFile(downloadURL, starttime)
}

func (ab *AudioBot) PlayAudioFile(url string, starttime int) {
	if ab.audio != nil {
		ab.audio.Stop()
		ab.audio = nil
	}
	audio, err := NewAudio(url, ab.VoiceConnection, ab.RoomChannel, starttime)
	if err != nil {
		fmt.Println("ERROR!!!!")
	}
	ab.audio = audio
	ab.audio.Play()
}
