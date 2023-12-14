package channels

import (
	"fmt"
	"math/rand"
	"sync"
	"watch2gether/pkg/channels/model"
	"watch2gether/pkg/channels/players"
	"watch2gether/pkg/media"

	"github.com/asdine/storm"
)

type Controller struct {
	Id       string
	players  map[players.PlayerType]players.Player
	previous []model.Event
	current  model.Event
	events   chan model.Event
	*storm.DB
	*sync.Mutex
}

func NewContoller(id string, store *storm.DB) *Controller {
	ctrl := Controller{
		Id:       id,
		DB:       store,
		players:  make(map[players.PlayerType]players.Player),
		previous: []model.Event{},
		Mutex:    &sync.Mutex{},
		events:   make(chan model.Event),
	}
	if _, err := ctrl.GetState(); err != nil {
		store.Save(model.NewPlayerState(id))
	}
	go ctrl.process()
	return &ctrl
}

func (controller *Controller) RemovePlayer(pType players.PlayerType) {
	if player, found := controller.players[pType]; found {
		player.Quit()
		delete(controller.players, pType)
	}
}

func (controller *Controller) GetPlayer(pType players.PlayerType) (players.Player, error) {
	if player, found := controller.players[pType]; found {
		return player, nil
	}
	return nil, fmt.Errorf("player not found")
}

func (controller *Controller) WithPlayer(pType players.PlayerType, player players.Player) *Controller {
	player.Start(controller.events)
	controller.players[pType] = player
	return controller
}

func (controller *Controller) GetState() (*model.PlayerState, error) {
	var state model.PlayerState
	err := controller.One("Id", controller.Id, &state)
	return &state, err
}

func (controller *Controller) update(state *model.PlayerState) error {
	return controller.Save(state)
}

func (controller *Controller) action(evt model.Event) {
	controller.Lock()
	defer controller.Unlock()
	fmt.Println("Processing Event:" + evt.Action)
	controller.previous = append(controller.previous, controller.current)
	controller.current = evt
}

func (c *Controller) Play() error {
	state, _ := c.GetState()
	if state.Current.Url == "" {
		state.Next()
	}
	c.action(c.current.WithAction(model.PLAY))
	for _, player := range c.players {
		go player.Play(c.current)
	}
	state.IsPlaying = true
	c.update(state)
	return c.Notify(c.current)
}

func (c *Controller) Pause() error {
	state, _ := c.GetState()
	c.action(c.current.WithAction(model.PAUSE))
	for _, player := range c.players {
		go player.Pause()
	}
	state.IsPlaying = false
	c.update(state)
	return c.Notify(c.current)
}

func (c *Controller) Stop() error {
	state, _ := c.GetState()
	c.action(c.current.WithAction(model.STOP))
	for _, player := range c.players {
		go player.Pause()
	}
	state.IsPlaying = false
	c.update(state)
	return c.Notify(c.current)
}

func (c *Controller) Skip() error {
	state, _ := c.GetState()
	c.action(c.current.WithAction(model.FINISHED))
	for _, player := range c.players {
		go player.Stop()
	}
	state.Next()
	c.update(state)
	c.Notify(c.current.WithAction(model.UPDATEQUEUE))
	return c.Play()
}

func (c *Controller) hasDoneEvent(evt model.Event) bool {
	for _, e := range c.previous {
		if e.ID == evt.ID {
			return true
		}
	}
	return c.current.ID == evt.ID && c.current.Action == evt.Action
}

func (c *Controller) process() {
	for evt := range c.events {
		if c.hasDoneEvent(evt) {
			continue
		}
		c.action(evt)
		switch evt.Action {
		case model.PLAY:
			c.Play()
		case model.PAUSE:
			c.Pause()
		case model.STOP:
			c.Stop()
		case model.FINISHED:
			c.Skip()
		}
	}
}

//QUECMD

func (controller *Controller) Add(meida []media.Media) error {
	controller.Lock()
	defer controller.Unlock()
	state, err := controller.GetState()
	if err != nil {
		return err
	}
	state.Queue = append(state.Queue, meida...)
	controller.update(state)
	return controller.Notify(controller.current.WithAction(model.UPDATEQUEUE))
}

func (controller *Controller) Move(srcIndex int, dstIndex int) error {
	state, err := controller.GetState()
	if err != nil {
		return err
	}
	state.Move(srcIndex, dstIndex)
	controller.update(state)
	return controller.Notify(controller.current.WithAction(model.UPDATEQUEUE))
}

func (controller *Controller) Remove(srcIndex int) error {
	state, err := controller.GetState()
	if err != nil {
		return err
	}
	state.Remove(srcIndex)
	controller.update(state)
	return controller.Notify(controller.current.WithAction(model.UPDATEQUEUE))
}

func (controller *Controller) UpdateQueue(meida []media.Media) error {
	controller.Lock()
	defer controller.Unlock()
	state, err := controller.GetState()
	if err != nil {
		return err
	}
	state.Queue = meida
	controller.update(state)
	return controller.Notify(controller.current.WithAction(model.UPDATEQUEUE))
}

func (controller *Controller) GetQueue() ([]media.Media, error) {
	controller.Lock()
	defer controller.Unlock()
	state, err := controller.GetState()
	if err != nil {
		return []media.Media{}, err
	}
	return state.Queue, nil
}

func (controller *Controller) GetCurrentVideo() (media.Media, error) {
	state, err := controller.GetState()
	if err != nil {
		return media.Media{}, err
	}
	return state.Current, nil
}

func (controller *Controller) Clear() error {
	controller.Lock()
	defer controller.Unlock()

	state, err := controller.GetState()
	if err != nil {
		return err
	}
	state.Queue = []media.Media{}
	controller.update(state)
	return controller.Notify(controller.current.WithAction(model.UPDATEQUEUE))
}

func (controller *Controller) SetLoop(loop bool) error {
	controller.Lock()
	defer controller.Unlock()
	state, err := controller.GetState()
	if err != nil {
		return err
	}
	state.Loop = loop
	return controller.update(state)
}

func (controller *Controller) Shuffle() error {
	controller.Lock()
	defer controller.Unlock()

	state, err := controller.GetState()
	if err != nil {
		return err
	}
	rand.Shuffle(len(state.Queue), func(i, j int) {
		state.Queue[i], state.Queue[j] = state.Queue[j], state.Queue[i]
	})
	controller.update(state)
	return controller.Notify(controller.current.WithAction(model.SHUFFLE))
}

func (controller *Controller) Notify(evt model.Event) error {
	for _, player := range controller.players {
		player.Notify(evt)
	}
	return nil
}
