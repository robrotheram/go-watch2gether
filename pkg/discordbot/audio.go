package discord

import (
	"fmt"
	"io"
	"log"
	"time"

	"github.com/jonas747/dca"
	"github.com/kkdai/youtube/v2"
)

// PlayAudioFile will play the given filename to the already connected
// Discord voice server/channel.  voice websocket and udp socket
// must already be setup before this will work.

func (db *DiscordBot) GetUserVoiceChannel(user string) (string, error) {
	for _, g := range db.session.State.Guilds {
		for _, v := range g.VoiceStates {
			if v.UserID == user {
				return v.ChannelID, nil
			}
		}
	}
	return "", fmt.Errorf("Channel Not found")
}

func (db *DiscordBot) PlayYoutube(videoURL string) {
	client := youtube.Client{}
	video, _ := client.GetVideo(videoURL)
	downloadURL, _ := client.GetStreamURL(video, &video.Formats[0])
	db.PlayAudioFile(downloadURL)
}

func (db *DiscordBot) PlayAudioFile(filename string) {
	// Send "speaking" packet over the voice websocket
	err := db.voice.Speaking(true)
	if err != nil {
		log.Fatal("Failed setting speaking", err)
	}

	// Send not "speaking" packet over the websocket when we finish
	defer db.voice.Speaking(false)

	opts := dca.StdEncodeOptions
	opts.RawOutput = true
	opts.Bitrate = 120

	encodeSession, err := dca.EncodeFile(filename, opts)
	if err != nil {
		log.Fatal("Failed creating an encoding session: ", err)
	}

	done := make(chan error)
	stream := dca.NewStream(encodeSession, db.voice, done)

	ticker := time.NewTicker(time.Second)

	for {
		select {
		case err := <-done:
			if err != nil && err != io.EOF {
				log.Fatal("An error occured", err)
			}

			// Clean up incase something happened and ffmpeg is still running
			encodeSession.Truncate()
			return
		case <-ticker.C:
			stats := encodeSession.Stats()
			playbackPosition := stream.PlaybackPosition()

			fmt.Printf("Playback: %10s, Transcode Stats: Time: %5s, Size: %5dkB, Bitrate: %6.2fkB, Speed: %5.1fx\r", playbackPosition, stats.Duration.String(), stats.Size, stats.Bitrate, stats.Speed)
		}
	}
}
