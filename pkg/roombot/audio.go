package roombot

import (
	"fmt"
	"io"
	"log"
	"sync"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/jonas747/dca"
)

type Audio struct {
	done     chan error
	stream   *dca.StreamingSession
	voice    *discordgo.VoiceConnection
	wg       sync.WaitGroup
	session  *dca.EncodeSession
	progress time.Duration
	url      string
	sync.Mutex
}

func NewAudio(url string, voice *discordgo.VoiceConnection) (*Audio, error) {
	opts := dca.StdEncodeOptions
	opts.RawOutput = true
	opts.Bitrate = 120

	encodeSession, err := dca.EncodeFile(url, opts)
	if err != nil {
		return nil, fmt.Errorf("Failed creating an encoding session: ", err)
	}
	return &Audio{
		url:     url,
		session: encodeSession,
		done:    make(chan error),
		voice:   voice,
	}, nil

}

func (audio *Audio) Unpause() {
	if audio.stream == nil {
		return
	}
	if audio.stream.Paused() {
		audio.voice.Speaking(true)
		audio.stream.SetPaused(true)
	}
}

func (audio *Audio) Paused() {
	if audio.stream == nil {
		return
	}
	if !audio.stream.Paused() {
		audio.voice.Speaking(false)
		audio.stream.SetPaused(true)
	}
}

func (audio *Audio) Play() {
	go func() {
		audio.voice.Speaking(true)
		defer audio.voice.Speaking(false)
		audio.PlayStream()
	}()
}

func (audio *Audio) Stop() {
	audio.session.Stop()
	<-audio.done
}

func (audio *Audio) PlayStream() {
	audio.stream = dca.NewStream(audio.session, audio.voice, audio.done)
	ticker := time.NewTicker(time.Second)
	for {
		select {
		case err := <-audio.done:
			if err != nil && err != io.EOF {
				log.Fatal("An error occured", err)
			}
			// Clean up incase something happened and ffmpeg is still running
			audio.session.Truncate()
			return
		case <-ticker.C:
			stats := audio.session.Stats()
			audio.progress = audio.stream.PlaybackPosition()
			fmt.Printf("Playback: %10s, Transcode Stats: Time: %5s, Size: %5dkB, Bitrate: %6.2fkB, Speed: %5.1fx\r", audio.progress, stats.Duration.String(), stats.Size, stats.Bitrate, stats.Speed)
		}
	}
}
