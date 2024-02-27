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
	id        string
	done      chan error
	stream    *dca.StreamingSession
	session   *dca.EncodeSession
	voice     *discordgo.VoiceConnection
	progress  media.MediaDuration
	running   bool
	seekTo    time.Duration
	startTime int
	exitcode  controllers.PlayerExitCode
}

func NewDiscordPlayer(id string, voice *discordgo.VoiceConnection) *DiscordPlayer {
	audio := &DiscordPlayer{
		id:       id,
		done:     make(chan error),
		voice:    voice,
		exitcode: controllers.STOP_EXITCODE,
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

func (player *DiscordPlayer) Id() string {
	return player.id
}

func (player *DiscordPlayer) Status() bool {
	return player.running
}

func (player *DiscordPlayer) Close() {
	player.exitcode = controllers.EXIT_EXITCODE
	player.Stop()
	player.voice.Disconnect()
}

func (player *DiscordPlayer) Seek(seconds time.Duration) {
	player.seekTo = seconds
	if player.session == nil {
		return
	}
	player.session.Stop()
	player.Finish()
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

func (player *DiscordPlayer) Play(url string, startTime int) (controllers.PlayerExitCode, error) {
	player.seekTo = -1
	if player.running {
		return controllers.STOP_EXITCODE, fmt.Errorf("playing already started")
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
		return controllers.STOP_EXITCODE, fmt.Errorf("failed creating an encoding session: %v", err)
	}
	player.session = encodeSession
	player.voice.Speaking(true)
	player.playStream()

	if player.seekTo > -1 {
		player.Finish()
		return player.Play(url, int(player.seekTo.Seconds()))
	}
	player.Finish()
	return player.exitcode, nil
}
