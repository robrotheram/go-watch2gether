package channels

import (
	"context"
	"log"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/robrotheram/dca"
)

type DiscordPlayer struct {
	voiceConnection *discordgo.VoiceConnection
	discord         *discordgo.Session
	stream          *dca.StreamingSession
	encodingSession *dca.EncodeSession
	*Base
}

func NewDiscordPlayer(id string, vc *discordgo.VoiceConnection, s *discordgo.Session) *DiscordPlayer {
	player := DiscordPlayer{
		voiceConnection: vc,
		discord:         s,
		Base:            NewBase(id),
	}
	return &player
}

func (dp *DiscordPlayer) check() {
	for _, g := range dp.discord.State.Guilds {
		if g.ID == dp.id {
			if len(g.VoiceStates) <= 1 {
				dp.Done()
			}
		}
	}
}

func (dp *DiscordPlayer) process() {
	for {
		dp.Lock()
		state, err := dp.GetState()
		if err != nil {
			return
		}
		dp.Unlock()
		if len(state.Current.Url) > 0 && state.State == PLAYING {
			log.Println(dp.Stream())
			if !state.Loop {
				dp.Next()
			}
		} else if len(state.Queue) > 0 && state.State == PLAYING {
			dp.Next()
		} else {
			time.Sleep(1 * time.Second)
		}
		dp.check()
	}
}
func (dp *DiscordPlayer) Run() error {
	go dp.process()
	return nil
}

func (dp *DiscordPlayer) Skip() error {
	if dp.encodingSession != nil {
		dp.encodingSession.Cleanup()
	}
	return nil
}

func (dp *DiscordPlayer) Play() error {
	state, err := dp.GetState()
	if err != nil {
		return err
	}
	if state.State == PAUSED {
		if dp.stream != nil {
			dp.stream.SetPaused(false)
		}
	} else if len(state.Current.AudioUrl) == 0 && len(state.Queue) > 0 {
		dp.Next()
	}

	state.State = PLAYING
	return dp.update(state)
}

func (dp *DiscordPlayer) Pause() error {
	dp.Lock()
	defer dp.Unlock()
	dp.stream.SetPaused(true)
	state, err := dp.GetState()
	if err != nil {
		return err
	}
	state.State = PAUSED
	return dp.update(state)
}

func (dp *DiscordPlayer) Stop() error {
	dp.Lock()
	defer dp.Unlock()
	if dp.encodingSession != nil {
		dp.encodingSession.Cleanup()
	}
	state, err := dp.GetState()
	if err != nil {
		return err
	}
	state.State = STOPPED
	return dp.update(state)
}

func (dp *DiscordPlayer) Done() error {
	dp.Stop()
	return dp.voiceConnection.Disconnect()
}

func (dp *DiscordPlayer) Duration() time.Duration {
	return dp.stream.PlaybackPosition()
}

func (dp *DiscordPlayer) updateDuration(ctx context.Context) {
	ticker := time.NewTicker(500 * time.Millisecond)
	for {
		select {
		case <-ticker.C:
			dp.Lock()
			state, err := dp.GetState()
			if err != nil {
				break
			}
			state.Proccessing = dp.Duration()
			dp.update(state)
			dp.Unlock()
		case <-ctx.Done():
			return
		}
	}
}

func (dp *DiscordPlayer) Stream() error {
	// Change these accordingly
	log.Println("Music Started")
	options := dca.StdEncodeOptions
	options.RawOutput = true
	options.Bitrate = 96
	options.Application = "lowdelay"

	var err error
	state, err := dp.GetState()
	if err != nil {
		return err
	}
	state.MediaRefresh()
	dp.encodingSession, err = dca.EncodeFile(state.Current.AudioUrl, options)
	if err != nil {
		return err
	}
	defer dp.encodingSession.Cleanup()

	dp.voiceConnection.Speaking(true)
	defer dp.voiceConnection.Speaking(false)

	done := make(chan error)
	ctx, cancel := context.WithCancel(context.Background())
	dp.stream = dca.NewStream(dp.encodingSession, dp.voiceConnection, done)
	go dp.updateDuration(ctx)
	err = <-done
	cancel()
	if err != nil {
		return err
	}
	log.Println("Music Ended")
	return nil
}
