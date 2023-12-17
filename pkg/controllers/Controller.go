package controllers

import (
	"fmt"
	"w2g/pkg/media"
)

type Controller struct {
	state         PlayerState
	player        Player
	notifications *Notify
	running       bool
}

func NewController(id string) *Controller {
	tracks := []media.Media{}
	return &Controller{
		state: PlayerState{
			Id:    id,
			Queue: tracks,
			State: STOP,
			Loop:  false,
		},
		notifications: NewNotifications(),
	}
}

func (c *Controller) Start() {
	if c.player == nil {
		return
	}
	if !c.running {
		if c.state.Current.Url == "" {
			c.state.Next()
		}
		c.running = true
		go c.progress()
	} else {
		c.unpause()
	}
	c.state.ChangeState(PLAY)
	c.Notify(PLAY_ACTION)
}

func (c *Controller) Stop() {
	c.running = false
	if c.player == nil {
		return
	}
	c.player.Stop()
	c.state.ChangeState(STOP)
	c.Notify(STOP_ACTION)
}

func (c *Controller) Pause() {
	c.player.Pause()
	c.state.ChangeState(PAUSE)
	c.Notify(PAUSE_ACTION)
}

func (c *Controller) unpause() {
	c.player.Unpause()
}

func (c *Controller) Add(url string, user string) error {
	tracks, err := media.NewVideo(url, user)
	if err != nil {
		return err
	}
	c.state.Add(tracks)
	c.Notify(UPDATEQUEUE_ACTION)
	return nil
}

func (c *Controller) Skip() {
	if c.running {
		c.player.Stop()
		c.Notify(SKIP_ACTION)
	}
}

func (c *Controller) Shuffle() {
	c.state.Shuffle()
	c.Notify(SHUFFLE_ACTION)
}

func (c *Controller) Loop() {
	c.state.Repeat()
	c.Notify(LOOP_ACTION)
}

func (c *Controller) Join(player Player) {
	c.player = player
	c.Notify(PLAYER_ACTION)
}

func (c *Controller) Leave() {
	c.player.Close()
	c.player = nil
	c.Notify(SHUFFLE_ACTION)
}

func (c *Controller) progress() {
	defer c.Stop()
	for {
		if len(c.state.Current.Url) == 0 || !c.running {
			c.Stop()
			return
		}
		audio := c.state.Current.GetAudioUrl()
		fmt.Println("playing: " + audio)
		err := c.player.Play(audio, 0)
		if err != nil {
			fmt.Println(err)
		}
		if !c.state.Loop {
			c.state.Next()
			c.Notify(UPDATEQUEUE_ACTION)
		}
	}
}

func (c *Controller) State() PlayerState {
	c.state.Active = c.IsActive()
	return c.state
}

func (c *Controller) Update(state PlayerState) {
	c.state = state
	c.Notify(UPDATEQUEUE_ACTION)
}

func (c *Controller) UpdateQueue(videos []media.Media) {
	c.state.Queue = videos
	c.Notify(UPDATEQUEUE_ACTION)
}

func (c *Controller) IsActive() bool {
	return c.player != nil
}

func (c *Controller) AddListner(listener Listener) {
	c.notifications.listners = append(c.notifications.listners, listener)
}

func (c *Controller) RemoveListner(listener Listener) {
	index := 0
	for _, i := range c.notifications.listners {
		if i != listener {
			c.notifications.listners[index] = i
			index++
		}
	}
	c.notifications.listners = c.notifications.listners[:index]
}

func (c *Controller) Notify(action Action) {
	state := c.State()
	c.notifications.events <- Event{
		ID:     state.Id,
		Action: action,
		State:  state,
	}
}
