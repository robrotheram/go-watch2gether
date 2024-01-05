package players

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"time"
	"w2g/pkg/controllers"
	"w2g/pkg/media"

	"github.com/bwmarrin/discordgo"
	"github.com/robrotheram/dca"
)

const DISCORD = controllers.PlayerType("DISCORD")

type DiscordPlayer struct {
	done      chan error
	stream    *dca.StreamingSession
	session   *dca.EncodeSession
	voice     *discordgo.VoiceConnection
	progress  media.MediaDuration
	running   bool
	startTime int
}

func NewDiscordPlayer(voice *discordgo.VoiceConnection) *DiscordPlayer {
	audio := &DiscordPlayer{
		done:  make(chan error),
		voice: voice,
	}
	return audio
}
func (player *DiscordPlayer) ParseDuration(url string) error {
	cmd := exec.Command(
		"ffprobe",
		"-i", url,
		"-show_entries", "format=duration",
		"-v", "quiet",
		"-of", "json",
	)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return err
	}

	type durationData struct {
		Format struct {
			Duration string `json:"duration"`
		} `json:"format"`
	}
	var d = durationData{}
	err = json.Unmarshal(out, &d)
	if err != nil {
		return err
	}

	duration, err := time.ParseDuration(d.Format.Duration + "s")
	player.progress = media.MediaDuration{
		Duration: duration,
		Progress: 0,
	}
	if err != nil {
		return err
	}
	return nil
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
			player.progress.Progress = player.stream.PlaybackPosition()
		}
	}
}

func (player *DiscordPlayer) Type() controllers.PlayerType {
	return DISCORD
}

func (player *DiscordPlayer) Status() bool {
	return player.running
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
	player.session.Stop()
	player.Finish()
}

func (player *DiscordPlayer) Progress() media.MediaDuration {
	return player.progress
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
	player.ParseDuration(url)
	encodeSession, err := dca.EncodeFile(url, opts)
	if err != nil {
		return fmt.Errorf("failed creating an encoding session: %v", err)
	}
	player.session = encodeSession

	player.voice.Speaking(true)
	defer player.voice.Speaking(false)
	defer player.Stop()
	player.playStream()
	return nil
}
