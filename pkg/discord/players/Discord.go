package players

import (
	"fmt"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/jonas747/dca"
)

type DiscordPlayer struct {
	done      chan error
	stream    *dca.StreamingSession
	session   *dca.EncodeSession
	voice     *discordgo.VoiceConnection
	progress  time.Duration
	running   bool
	startTime int
}

func NewDiscordPlayer(voice *discordgo.VoiceConnection) *DiscordPlayer {
	audio := &DiscordPlayer{
		done:  make(chan error),
		voice: voice,
	}
	// audio.start()
	return audio
}

func (player *DiscordPlayer) playStream() {
	player.stream = dca.NewStream(player.session, player.voice, player.done)
	ticker := time.NewTicker(time.Second)
	for {
		select {
		case <-player.done:
			// Clean up incase something happened and ffmpeg is still running
			player.Finish()
			return
		case <-ticker.C:
			//stats := audio.session.Stats()
			player.progress = player.stream.PlaybackPosition()
			//fmt.Printf("Playback: %10s, Transcode Stats: Time: %5s, Size: %5dkB, Bitrate: %6.2fkB, Speed: %5.1fx\r", audio.progress, stats.Duration.String(), stats.Size, stats.Bitrate, stats.Speed)
			// progress := float64(0)
			// if player.Duration > 0 && player.progress > 0 {
			// 	progress = (float64(player.startTime) + audio.progress.Seconds()) / audio.Duration.Seconds()
			// }
			//Only send event if we are playing. Fixes issue if ticker fires after audio.done
			// if player.Playing {
			// 	player.bot.SendToRoom(CreateBotUpdateEvent(media.Seek{
			// 		ProgressPct: progress,
			// 		ProgressSec: float64(audio.startTime) + audio.progress.Seconds(),
			// 	}))
			// }
		}
	}
}

func (player *DiscordPlayer) Close() {
	player.Stop()
	player.voice.Disconnect()
}

func (player *DiscordPlayer) Finish() {
	player.session.Cleanup()
	player.session.Truncate()
	player.running = false
}

func (player *DiscordPlayer) Unpause() {
	if player.stream == nil {
		return
	}
	player.voice.Speaking(true)
	player.stream.SetPaused(false)
}

func (player *DiscordPlayer) Pause() {
	if player.stream == nil {
		return
	}
	player.voice.Speaking(false)
	player.stream.SetPaused(true)
}

func (player *DiscordPlayer) Stop() {
	if player.session == nil {
		return
	}
	err := player.session.Stop()
	if err != nil {
		fmt.Println(err)
	}
	player.Finish()
}

func (player *DiscordPlayer) Play(url string, startTime int) error {
	if player.running {
		return fmt.Errorf("playing already started")
	}
	opts := dca.StdEncodeOptions
	opts.RawOutput = true
	opts.Bitrate = 96
	opts.StartTime = startTime
	opts.Application = "audio"
	opts.PacketLoss = 10
	player.startTime = startTime

	encodeSession, err := dca.EncodeFile(url, opts)
	if err != nil {
		return fmt.Errorf("failed creating an encoding session: %v", err)
	}
	player.session = encodeSession

	player.voice.Speaking(true)
	defer player.voice.Speaking(false)
	player.playStream()
	return nil
}
