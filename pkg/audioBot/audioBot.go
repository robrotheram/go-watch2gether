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

func (ab *AudioBot) RegisterToRoom(rc chan []byte) {
	ab.RoomChannel = rc
	ab.audio = NewAudio(rc, ab.VoiceConnection)
}

func (ab *AudioBot) Send(evt events.Event) {
	go func() { ab.handleEvent(evt) }()
}

func (ab *AudioBot) sendToChannel(msg string) {
	ab.session.ChannelMessageSend(ab.notficationChannelID, msg)
}

func (ab *AudioBot) handleEvent(evt events.Event) {
	if ab.audio == nil {
		return
	}
	switch evt.Action {
	case events.EVNT_UPDATE_QUEUE:
		ab.sendToChannel(fmt.Sprintf("Queue Updated by: %s", evt.Watcher.Username))

	case events.EVT_VIDEO_CHANGE:
		ab.PlayAudio(evt.CurrentVideo, 0)

	case events.EVNT_PLAYING:
		if !ab.audio.Playing {
			ab.PlayAudio(evt.CurrentVideo, 0)
		} else {
			ab.audio.Unpause()
		}
		ab.sendToChannel(fmt.Sprintf("User: %s Started the video", evt.Watcher.Username))
	case events.EVNT_PAUSING:
		ab.audio.Paused()
		ab.sendToChannel(fmt.Sprintf("User: %s Paused the video", evt.Watcher.Username))

	case events.EVNT_SEEK_TO_USER:
		ab.audio.Stop()
		ab.PlayAudio(evt.CurrentVideo, int(evt.Seek.ProgressSec))
		ab.sendToChannel(fmt.Sprintf("User: %s Seeked Video to %f", evt.Watcher.Username, evt.Seek.ProgressPct*100))

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
	if ab.audio == nil {
		fmt.Println("Bot not connected to Room")
	}
	if ab.audio.Playing {
		ab.audio.Stop()
		fmt.Println("Stopping Audio")
	}
	err := ab.audio.Play(url, starttime)
	if err != nil {
		fmt.Printf("Error Encoding URL : %v \n", err)
		return
	}
	err = ab.audio.Start()
	if err != nil {
		fmt.Printf("Error Starting Audio : %v \n", err)
		return
	}
}

func (ab *AudioBot) Disconnect() error {
	if ab.audio != nil {
		ab.audio.Stop()
	}
	return ab.VoiceConnection.Disconnect()
}
