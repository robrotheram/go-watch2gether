package players

import (
	"fmt"
	"sync"
	"time"

	"watch2gether/pkg/channels/model"

	"github.com/bwmarrin/discordgo"
	"github.com/robrotheram/dca"
	log "github.com/sirupsen/logrus"
)

var DISCORD = PlayerType("discord")

type DiscordPlayer struct {
	voiceConnection *discordgo.VoiceConnection
	discord         *discordgo.Session
	stream          *dca.StreamingSession
	encodingSession *dca.EncodeSession
	notify          chan model.Event
	event           model.Event
	*sync.Mutex
}

func NewDiscordPlayer(vc *discordgo.VoiceConnection, s *discordgo.Session) *DiscordPlayer {
	player := DiscordPlayer{
		voiceConnection: vc,
		discord:         s,
		Mutex:           &sync.Mutex{},
	}
	return &player
}

func (dp *DiscordPlayer) Start(notify chan model.Event) {
	dp.notify = notify
}

func (dp *DiscordPlayer) Notify(event model.Event) {
	// interaction := discordgo.InteractionResponse{
	// 	Type: discordgo.InteractionResponseChannelMessageWithSource,
	// 	Data: &discordgo.InteractionResponseData{
	// 		Content: "Hello",
	// 	},
	// }
	// fmt.Println("HHHHHHHHHHHHHHHHHHHHHHHHHHHHHH")
	// err := dp.discord.InteractionRespond(&discordgo.Interaction{
	// 	GuildID: "",
	// 	ChannelID: "",

	// }, &interaction)
	//dp.discord.ChannelMessageSend("368777776536223756", "HELLO")
}

func (dp *DiscordPlayer) Skip() error {
	if dp.encodingSession != nil {
		dp.encodingSession.Cleanup()
	}
	return nil
}

func (dp *DiscordPlayer) Pause() error {
	if dp.stream == nil {
		return nil
	}
	dp.Lock()
	defer dp.Unlock()
	dp.stream.SetPaused(true)
	return nil
}

func (dp *DiscordPlayer) Stop() error {
	dp.Lock()
	defer dp.Unlock()
	if dp.encodingSession != nil {
		dp.encodingSession.Cleanup()
	}
	return nil
}

func (dp *DiscordPlayer) Quit() error {
	dp.Stop()
	return dp.voiceConnection.Disconnect()
}

func (dp *DiscordPlayer) Duration() time.Duration {
	if dp.encodingSession != nil {
		pos := dp.stream.PlaybackPosition()
		log.Infof("DS TIME Left %d", dp.event.Duration-pos)
		return pos
	}
	return time.Duration(0)
}

func (dp *DiscordPlayer) onFinish() {
	dp.notify <- dp.event.WithAction(model.FINISHED)
}

func (dp *DiscordPlayer) Play(event model.Event) error {
	dp.event = event
	video := event.Media
	if dp.stream != nil {
		if dp.stream.Paused() {
			dp.stream.SetPaused(false)
			return nil
		}
		if dp.encodingSession != nil {
			dp.encodingSession.Cleanup()
		}
	}
	err := video.Refresh()
	if err != nil {
		return err
	}
	log.Println("Music Started: " + video.AudioUrl)
	options := dca.StdEncodeOptions
	options.RawOutput = true
	options.Bitrate = 96
	options.Application = "lowdelay"

	dp.encodingSession, err = dca.EncodeFile(video.AudioUrl, options)
	if err != nil {
		return err
	}
	defer dp.encodingSession.Cleanup()

	dp.voiceConnection.Speaking(true)
	defer dp.voiceConnection.Speaking(false)

	done := make(chan error)
	dp.stream = dca.NewStream(dp.encodingSession, dp.voiceConnection, done)
	//Wait for stream to finish
	err = <-done
	fmt.Println("Stream ERROR: " + err.Error())
	log.Println("Music Ended")
	dp.onFinish()
	return nil
}
