package controllers

import (
	"fmt"
	"sync"
	"w2g/pkg/media"
)

type Controller struct {
	state   PlayerState
	player  Player
	running bool
	sync.Mutex
}

func NewController() *Controller {
	tracks := []media.Media{}
	return &Controller{
		state: PlayerState{
			Queue: tracks,
		},
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
}

func (c *Controller) Stop() {
	defer c.Unlock()
	c.Lock()
	c.running = false
	c.player.Stop()
}

func (c *Controller) Pause() {
	c.player.Pause()
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
	return nil
}

func (c *Controller) Skip() {
	if c.running {
		c.player.Stop()
	}
}

func (c *Controller) Join(player Player) {
	c.player = player
}

func (c *Controller) Leave() {
	c.player.Close()
}

func (c *Controller) progress() {
	defer c.Stop()
	for {
		if len(c.state.Current.Url) == 0 || !c.running {
			c.Stop()
			return
		}
		fmt.Println("playing: " + c.state.Current.Url)
		audio := c.state.Current.GetAudioUrl()
		err := c.player.Play(audio, 0)
		if err != nil {
			fmt.Println(err)
		}
		c.state.Next()
	}
}

func (c *Controller) State() PlayerState {
	return c.state
}

func (c *Controller) Update(state PlayerState) {
	c.state = state
}

func (c *Controller) IsActive() bool {
	return c.player != nil
}
