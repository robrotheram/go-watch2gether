package players

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
	voiceConnection *discordgo.VoiceConnection
	stream          *dca.StreamingSession
	encodingSession *dca.EncodeSession
	store           *storm.DB
	state           *Player
	*sync.Mutex
}

func NewDiscordPlayer(id string, vc *discordgo.VoiceConnection) *DiscordPlayer {
	player := DiscordPlayer{
		voiceConnection: vc,
		state: &Player{
			Id:    id,
			State: STOPPED,
			Queue: []media.Media{},
		},
		Mutex: &sync.Mutex{},
	}
	go player.process()
	return &player
}

func (dp *DiscordPlayer) process() {
	for {
		dp.Lock()
		state := dp.state
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
	}
}

func (dp *DiscordPlayer) SetStore(store *storm.DB) {
	dp.store = store
	dp.Load()
	dp.update()
}

func (dp *DiscordPlayer) update() {
	log.Println(dp.store.Save(dp.state))
}

func (dp *DiscordPlayer) Load() {
	var p Player
	err := dp.store.One("Id", dp.state.Id, &p)
	if err == nil {
		dp.state = &p
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
	dp.state.Queue = videos
	dp.update()
}

func (dp *DiscordPlayer) SetLoop(loop bool) {
	dp.Lock()
	defer dp.Unlock()
	dp.state.Loop = loop
	dp.update()
}

func (dp *DiscordPlayer) UpdaetState(state *Player) {
	dp.Lock()
	defer dp.Unlock()
	dp.state = state
	dp.update()
}

func (dp *DiscordPlayer) Play() {
	if dp.state.State == PAUSED {
		if dp.stream != nil {
			dp.stream.SetPaused(false)
		}
	} else if len(dp.state.Current.AudioUrl) == 0 && len(dp.state.Queue) > 0 {
		dp.Next()
	}

	dp.state.State = PLAYING
	dp.update()
}

func (dp *DiscordPlayer) Add(meida []media.Media) {
	dp.Lock()
	defer dp.Unlock()
	dp.state.Queue = append(dp.state.Queue, meida...)
	dp.update()
}

func (dp *DiscordPlayer) Next() {
	dp.Lock()
	defer dp.Unlock()
	dp.state.Current = media.Media{}
	if len(dp.state.Queue) > 0 {
		dp.state.Current, dp.state.Queue = dp.state.Queue[0], dp.state.Queue[1:]
	}
	dp.update()
}

func (dp *DiscordPlayer) Shuffle() {
	dp.Lock()
	defer dp.Unlock()
	rand.Shuffle(len(dp.state.Queue), func(i, j int) {
		dp.state.Queue[i], dp.state.Queue[j] = dp.state.Queue[j], dp.state.Queue[i]
	})
	dp.update()

}

func (dp *DiscordPlayer) GetState() *Player {
	dp.Lock()
	defer dp.Unlock()
	return dp.state
}

func (dp *DiscordPlayer) GetQueue() []media.Media {
	dp.Lock()
	defer dp.Unlock()
	return dp.state.Queue
}

func (dp *DiscordPlayer) GetCurrentVideo() media.Media {
	return dp.state.Current
}

func (dp *DiscordPlayer) Clear() {
	dp.Lock()
	defer dp.Unlock()
	dp.state.Queue = []media.Media{}
	dp.update()
}

func (dp *DiscordPlayer) Pause() {
	dp.Lock()
	defer dp.Unlock()
	dp.stream.SetPaused(true)
	dp.state.State = PAUSED
	dp.update()
}

func (dp *DiscordPlayer) Stop() {
	dp.Lock()
	defer dp.Unlock()
	if dp.encodingSession != nil {
		dp.encodingSession.Cleanup()
	}
	dp.state.State = STOPPED
	dp.update()
}

func (dp *DiscordPlayer) Done() {
	dp.Stop()
	dp.voiceConnection.Disconnect()
	dp.update()
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
			dp.state.Proccessing = dp.Duration()
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
	dp.state.MediaRefresh()
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
