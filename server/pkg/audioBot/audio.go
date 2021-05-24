package audioBot

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"sync"
	"time"
	"watch2gether/pkg/media"

	"github.com/bwmarrin/discordgo"
	"github.com/robrotheram/dca"
	log "github.com/sirupsen/logrus"
)

type Audio struct {
	done      chan error
	stream    *dca.StreamingSession
	voice     *discordgo.VoiceConnection
	wg        sync.WaitGroup
	session   *dca.EncodeSession
	progress  time.Duration
	Url       string
	Duration  time.Duration
	Playing   bool
	startTime int
	sync.Mutex
	bot *AudioBot
}

func ParseDuration(url string) (time.Duration, error) {
	cmd := exec.Command(
		"ffprobe",
		"-i", url,
		"-show_entries", "format=duration",
		"-v", "quiet",
		"-of", "json",
	)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return 0, err
	}

	type duration struct {
		Format struct {
			Duration string `json:"duration"`
		} `json:"format"`
	}
	var d = duration{}
	err = json.Unmarshal(out, &d)
	if err != nil {
		return 0, err
	}
	return time.ParseDuration(d.Format.Duration + "s")
}

func NewAudio(bot *AudioBot, voice *discordgo.VoiceConnection) *Audio {
	audio := Audio{
		done:  make(chan error),
		voice: voice,
		bot:   bot,
	}
	return &audio
}

func (audio *Audio) Play(url string, startTime int) error {
	if audio.Playing {
		return fmt.Errorf("Playing already started")
	}
	audio.Stop()
	audio.Lock()
	defer audio.Unlock()
	opts := dca.StdEncodeOptions
	opts.RawOutput = true
	opts.Bitrate = 96
	opts.StartTime = startTime
	opts.Application = "lowdelay"
	audio.startTime = startTime

	encodeSession, err := dca.EncodeFile(url, opts)
	if err != nil {
		return fmt.Errorf("Failed creating an encoding session: %v", err)
	}
	go func() {
		audio.Duration, _ = ParseDuration(audio.Url)
	}()
	audio.session = encodeSession
	audio.Url = url
	return nil
}

func (audio *Audio) Start() error {
	audio.Lock()
	start := !audio.Playing
	audio.Playing = true
	audio.Unlock()
	if start {
		audio.wg.Add(1)
		go func() {
			audio.voice.Speaking(true)
			defer audio.voice.Speaking(false)
			audio.PlayStream()
			audio.Lock()
			audio.Playing = false
			defer audio.wg.Done()
			audio.Unlock()
		}()
		return nil
	} else {
		return fmt.Errorf("Playing already started")
	}
}

func (audio *Audio) Unpause() {
	if audio.stream == nil {
		return
	}
	audio.voice.Speaking(true)
	audio.stream.SetPaused(false)
}

func (audio *Audio) Paused() {
	if audio.stream == nil {
		return
	}
	audio.voice.Speaking(false)
	audio.stream.SetPaused(true)
}

func (audio *Audio) Stop() {
	audio.Paused()
	if audio.session == nil {
		return
	}
	err := audio.session.Stop()
	if err != nil {
		log.Debugf("Error Stoppin Session %v", err)
	}
	audio.Lock()
	audio.Playing = false
	audio.session.Cleanup()
	audio.session.Truncate()
	audio.Unlock()
}

func (audio *Audio) PlayStream() {

	audio.stream = dca.NewStream(audio.session, audio.voice, audio.done)
	ticker := time.NewTicker(time.Second)
	for {
		if !audio.Playing {
			audio.session.Truncate()
			audio.bot.SendToRoom(CreateBotFinishEvent())
			audio.Playing = false
			return
		}
		select {
		case <-audio.done:
			// Clean up incase something happened and ffmpeg is still running
			audio.session.Truncate()
			audio.Playing = false
			return
		case <-ticker.C:
			//stats := audio.session.Stats()
			audio.progress = audio.stream.PlaybackPosition()
			//fmt.Printf("Playback: %10s, Transcode Stats: Time: %5s, Size: %5dkB, Bitrate: %6.2fkB, Speed: %5.1fx\r", audio.progress, stats.Duration.String(), stats.Size, stats.Bitrate, stats.Speed)
			progress := float64(0)
			if audio.Duration > 0 && audio.progress > 0 {
				progress = (float64(audio.startTime) + audio.progress.Seconds()) / audio.Duration.Seconds()
			}
			audio.bot.SendToRoom(CreateBotUpdateEvent(media.Seek{
				ProgressPct: progress,
				ProgressSec: float64(audio.startTime) + audio.progress.Seconds(),
			}))
		}
	}
}
