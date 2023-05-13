package channels

import (
	"context"
	"log"
	"math/rand"
	"sync"
	"time"
	"watch2gether/pkg/media"

	"github.com/asdine/storm"
	"github.com/bwmarrin/discordgo"
	"github.com/robrotheram/dca"
)

type DiscordPlayer struct {
	id              string
	voiceConnection *discordgo.VoiceConnection
	discord         *discordgo.Session
	stream          *dca.StreamingSession
	encodingSession *dca.EncodeSession
	*storm.DB
	*sync.Mutex
}

func NewDiscordPlayer(id string, vc *discordgo.VoiceConnection, s *discordgo.Session) *DiscordPlayer {
	player := DiscordPlayer{
		id:              id,
		voiceConnection: vc,
		Mutex:           &sync.Mutex{},
	}
	go player.process()
	return &player
}

func (dp *DiscordPlayer) check() {
	for _, g := range dp.discord.State.Guilds {
		if g.ID == dp.id {
			if len(g.VoiceStates) <= 1 {
				dp.disconnect()
			}
		}
	}
}

func (dp *DiscordPlayer) disconnect() {
	dp.Stop()
	dp.voiceConnection.Disconnect()
}

func (dp *DiscordPlayer) update(state *Player) {
	dp.Save(state)
}

func (dp *DiscordPlayer) SetStore(db *storm.DB) {
	dp.DB = db
}

func (dp *DiscordPlayer) GetState() *Player {
	var state *Player
	dp.One("Id", dp.id, state)
	return state
}

func (dp *DiscordPlayer) process() {
	for {
		dp.Lock()
		state := dp.GetState()
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

func (dp *DiscordPlayer) Skip() {
	if dp.encodingSession != nil {
		dp.encodingSession.Cleanup()
	}
}

func (dp *DiscordPlayer) UpdateQueue(videos []media.Media) {
	dp.Lock()
	defer dp.Unlock()

	state := dp.GetState()
	state.Queue = videos
	dp.update(state)
}

func (dp *DiscordPlayer) SetLoop(loop bool) {
	dp.Lock()
	defer dp.Unlock()
	state := dp.GetState()
	state.Loop = loop
	dp.update(state)
}

func (dp *DiscordPlayer) Play() {
	state := dp.GetState()
	if state.State == PAUSED {
		if dp.stream != nil {
			dp.stream.SetPaused(false)
		}
	} else if len(state.Current.AudioUrl) == 0 && len(state.Queue) > 0 {
		dp.Next()
	}

	state.State = PLAYING
	dp.update(state)
}

func (dp *DiscordPlayer) Add(meida []media.Media) {
	dp.Lock()
	defer dp.Unlock()
	state := dp.GetState()
	state.Queue = append(state.Queue, meida...)
	dp.update(state)
}

func (dp *DiscordPlayer) Next() {
	dp.Lock()
	defer dp.Unlock()

	state := dp.GetState()
	state.Current = media.Media{}
	if len(state.Queue) > 0 {
		state.Current, state.Queue = state.Queue[0], state.Queue[1:]
	}
	dp.update(state)
}

func (dp *DiscordPlayer) Shuffle() {
	dp.Lock()
	defer dp.Unlock()

	state := dp.GetState()
	rand.Shuffle(len(state.Queue), func(i, j int) {
		state.Queue[i], state.Queue[j] = state.Queue[j], state.Queue[i]
	})
	dp.update(state)

}

func (dp *DiscordPlayer) GetQueue() []media.Media {
	dp.Lock()
	defer dp.Unlock()
	return dp.GetState().Queue
}

func (dp *DiscordPlayer) GetCurrentVideo() media.Media {
	return dp.GetState().Current
}

func (dp *DiscordPlayer) Clear() {
	dp.Lock()
	defer dp.Unlock()

	state := dp.GetState()
	state.Queue = []media.Media{}
	dp.update(state)
}

func (dp *DiscordPlayer) Pause() {
	dp.Lock()
	defer dp.Unlock()
	dp.stream.SetPaused(true)
	state := dp.GetState()
	state.State = PAUSED
	dp.update(state)
}

func (dp *DiscordPlayer) Stop() {
	dp.Lock()
	defer dp.Unlock()
	if dp.encodingSession != nil {
		dp.encodingSession.Cleanup()
	}
	state := dp.GetState()
	state.State = STOPPED
	dp.update(state)
}

func (dp *DiscordPlayer) Done() {
	dp.Stop()
	dp.voiceConnection.Disconnect()
}

func (dp *DiscordPlayer) Duration() time.Duration {
	return dp.stream.PlaybackPosition()
}

func (dp *DiscordPlayer) Move(srcIndex int, dstIndex int) {
	state := dp.GetState()
	state.Move(srcIndex, dstIndex)
	dp.update(state)
}

func (dp *DiscordPlayer) Remove(srcIndex int) {
	state := dp.GetState()
	state.Remove(srcIndex)
	dp.update(state)
}

func (dp *DiscordPlayer) updateDuration(ctx context.Context) {
	ticker := time.NewTicker(500 * time.Millisecond)
	for {
		select {
		case <-ticker.C:
			dp.Lock()
			state := dp.GetState()
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
	state := dp.GetState()
	state.MediaRefresh()
	dp.encodingSession, err = dca.EncodeFile(dp.GetCurrentVideo().AudioUrl, options)
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
