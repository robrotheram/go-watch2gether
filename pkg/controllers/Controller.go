package controllers

import (
	"fmt"
	"time"
	"w2g/pkg/media"
)

const SYSTEM = "system"

type Controller struct {
	state         PlayerState
	player        Player
	notifications *Notify
	running       bool
}

func NewController(id string) *Controller {
	tracks := []media.Media{}
	contoller := Controller{
		state: PlayerState{
			Id:    id,
			Queue: tracks,
			State: STOP,
			Loop:  false,
		},
		notifications: NewNotifications(),
	}
	contoller.AddListner(&Auditing{})
	return &contoller
}

func (c *Controller) Start(user string) {
	if c.player == nil {
		return
	}
	if !c.running {
		if c.state.Current.Url == "" {
			c.state.Next()
		}
		c.running = true
		go c.progress()
		go c.duration()
	} else {
		c.unpause()
	}
	c.state.ChangeState(PLAY)
	c.Notify(PLAY_ACTION, user)
}

func (c *Controller) Stop(user string) {
	c.running = false
	if c.player == nil {
		return
	}
	c.player.Stop()
	c.state.ChangeState(STOP)
	c.Notify(STOP_ACTION, user)
}

func (c *Controller) Pause(user string) {
	c.player.Pause()
	c.state.ChangeState(PAUSE)
	c.Notify(PAUSE_ACTION, user)
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
	c.Notify(ADD_QUEUE, user)
	return nil
}

func (c *Controller) Skip(user string) {
	if c.running {
		c.player.Stop()
		c.Notify(SKIP_ACTION, user)
	}
}

func (c *Controller) Shuffle(user string) {
	c.state.Shuffle()
	c.Notify(SHUFFLE_ACTION, user)
}

func (c *Controller) Loop(user string) {
	c.state.Repeat()
	c.Notify(LOOP_ACTION, user)
}

func (c *Controller) Join(player Player, user string) {
	c.player = player
	c.Notify(PLAYER_ACTION, user)
}

func (c *Controller) Leave(user string) {
	c.player.Close()
	c.player = nil
	c.Notify(SHUFFLE_ACTION, user)
}

func (c *Controller) progress() {
	defer c.Stop(SYSTEM)
	for {
		if len(c.state.Current.Url) == 0 || !c.running {
			c.Stop(SYSTEM)
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
			c.Notify(UPDATE_QUEUE, SYSTEM)
		}
	}
}

func (c *Controller) duration() {
	ticker := time.NewTicker(1 * time.Second)
	for range ticker.C {
		if !c.running {
			ticker.Stop()
			return
		}
		c.state.Current.Progress = c.player.Progress()
		c.Notify(UPDATE_DURATION, SYSTEM)
	}
}

func (c *Controller) State() PlayerState {
	c.state.Active = c.IsActive()
	return c.state
}

func (c *Controller) Update(state PlayerState, user string) {
	c.state = state
	c.Notify(UPDATE, user)
}

func (c *Controller) UpdateQueue(videos []media.Media, user string) {
	c.state.Queue = videos
	c.Notify(UPDATE_QUEUE, user)
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

func (c *Controller) Notify(action ActionType, user string) {
	state := c.State()
	c.notifications.events <- Event{
		ID: state.Id,
		Action: Action{
			ActionType: action,
			User:       user,
		},
		State: state,
	}
}
