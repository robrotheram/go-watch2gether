package channels

import (
	"context"
	"time"

	"github.com/gorilla/websocket"
)

type WebPlayer struct {
	ws *websocket.Conn
	*Base
}

func NewWebPlayer(id string, ws *websocket.Conn) *WebPlayer {
	player := WebPlayer{
		ws:   ws,
		Base: NewBase(id),
	}
	return &player
}

func (wp *WebPlayer) check() {
}

func (wp *WebPlayer) process() {
	for {
		select {
		case <-wp.quit:
			return
		default:
			wp.Lock()

			// state, err := wp.GetState()
			// if err != nil {
			// 	return
			// }
			wp.Unlock()
			// if len(state.Current.Url) > 0 && state.State == PLAYING {
			// 	log.Println(wp.Stream())
			// 	if !state.Loop {
			// 		wp.Next()
			// 	}
			// } else if len(state.Queue) > 0 && state.State == PLAYING {
			// 	wp.Next()
			// } else {
			// 	time.Sleep(1 * time.Second)
			// }
			// wp.ws.WriteJSON(state)
			// wp.Broadcast("SOMETHING")
			time.Sleep(1 * time.Second)
			wp.check()
		}
	}
}
func (wp *WebPlayer) Run() error {
	go wp.process()
	return nil
}

func (wp *WebPlayer) Skip() error {
	return nil
}

func (wp *WebPlayer) Play() error {
	state, err := wp.GetState()
	if err != nil {
		return err
	}
	if state.State == PAUSED {
		// if wp.stream != nil {
		// 	wp.stream.SetPaused(false)
		// }
	} else if len(state.Current.AudioUrl) == 0 && len(state.Queue) > 0 {
		wp.Next()
	}

	state.State = PLAYING
	return wp.update(state)
}

func (wp *WebPlayer) Pause() error {
	wp.Lock()
	defer wp.Unlock()
	// wp.stream.SetPaused(true)
	state, err := wp.GetState()
	if err != nil {
		return err
	}
	state.State = PAUSED
	return wp.update(state)
}

func (wp *WebPlayer) Stop() error {
	wp.Lock()
	defer wp.Unlock()
	// if wp.encodingSession != nil {
	// 	wp.encodingSession.Cleanup()
	// }
	state, err := wp.GetState()
	if err != nil {
		return err
	}
	state.State = STOPPED
	return wp.update(state)
}

func (wp *WebPlayer) Done() error {
	wp.Stop()
	go func() { wp.quit <- true }()
	return wp.ws.Close()
}

func (wp *WebPlayer) Duration() time.Duration {
	return 0
}

func (wp *WebPlayer) updateDuration(ctx context.Context) {
	ticker := time.NewTicker(500 * time.Millisecond)
	for {
		select {
		case <-ticker.C:
			wp.Lock()
			state, err := wp.GetState()
			if err != nil {
				break
			}
			state.Proccessing = wp.Duration()
			wp.update(state)
			wp.Unlock()
		case <-ctx.Done():
			return
		}
	}
}

func (wp *WebPlayer) Stream() error {
	return nil
}

func (wp *WebPlayer) GetType() RoomType {
	return WEBSOCKET
}

func (wp *WebPlayer) Notify(evt EventType) {
	switch evt {
	case PLAY:
		wp.Play()
	case PAUSE:
		wp.Pause()
	case SKIP:
		wp.Skip()
	case STOP:
		wp.Stop()
	}
}
