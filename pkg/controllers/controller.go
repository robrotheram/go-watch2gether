package controllers

import (
	"time"
	"w2g/pkg/media"

	"github.com/asdine/storm"
	"github.com/google/uuid"
)

const SYSTEM = "system"

type Controller struct {
	state         PlayerState
	players       *Players
	notifications *Notify
	running       bool
	db            *storm.DB
}

func NewController(id string, store *storm.DB) *Controller {
	contoller := Controller{
		players:       newPlayers(),
		notifications: NewNotifications(),
		db:            store,
	}
	contoller.load(id)
	contoller.AddListner(uuid.NewString(), &Auditing{})
	return &contoller
}

func (c *Controller) save() error {
	return c.db.Save(&c.state)
}

func (c *Controller) load(id string) {
	var state PlayerState
	err := c.db.One("ID", id, &state)
	if err != nil {
		state = PlayerState{
			ID:    id,
			Queue: []media.Media{},
			State: STOP,
			Loop:  false,
		}
		c.save()
	}
	state.ChangeState(STOP)
	c.state = state
}

func (c *Controller) Start(user string) {
	if c.players.Empty() {
		return
	}
	if !c.running {
		if c.state.Current.Url == "" {
			c.state.Next()
		}
		c.running = true
		go c.progress()
		go c.duration()
	} else if c.state.State == PAUSE {
		c.players.Unpause()
	}
	c.state.ChangeState(PLAY)
	c.save()
	c.Notify(PLAY_ACTION, user)
}

func (c *Controller) Stop(user string) {
	c.running = false
	if c.players.Empty() {
		return
	}
	c.players.Stop()
	c.state.ChangeState(STOP)
	c.save()
	c.Notify(STOP_ACTION, user)
}

func (c *Controller) Pause(user string) {
	c.players.Pause()
	c.state.ChangeState(PAUSE)
	c.save()
	c.Notify(PAUSE_ACTION, user)
}

func (c *Controller) Add(url string, user string) error {
	tracks, err := media.NewVideo(url, user)
	if err != nil {
		return err
	}
	c.state.Add(tracks)
	c.save()
	c.Notify(ADD_QUEUE, user)
	return nil
}

func (c *Controller) Skip(user string) {
	if c.running {
		c.players.Stop()
		c.Notify(SKIP_ACTION, user)
	}
}

func (c *Controller) Shuffle(user string) {
	c.state.Shuffle()
	c.save()
	c.Notify(SHUFFLE_ACTION, user)
}

func (c *Controller) Loop(user string) {
	c.state.Repeat()
	c.save()
	c.Notify(LOOP_ACTION, user)
}

func (c *Controller) UpdateQueue(videos []media.Media, user string) {
	c.state.Queue = videos
	c.save()
	c.Notify(UPDATE_QUEUE, user)
}

func (c *Controller) State() PlayerState {
	c.state.Active = !c.players.Empty()
	return c.state
}

func (c *Controller) Update(state PlayerState, user string) {
	c.state = state
	c.save()
	c.Notify(UPDATE, user)
}

func (c *Controller) Join(player Player, user string) {
	c.players.Add(player)
	c.Notify(PLAYER_ACTION, user)
}

func (c *Controller) Leave(pType PlayerType, user string) {
	c.players.Remvoe(pType)
	c.Notify(LEAVE_ACTION, user)
}

func (c *Controller) ContainsPlayer(pType PlayerType) bool {
	if _, ok := c.players.players[pType]; ok {
		return true
	}
	return false
}

func (c *Controller) progress() {
	defer c.Stop(SYSTEM)
	for {
		if len(c.state.Current.Url) == 0 || !c.running || c.players.Empty() {
			c.Stop(SYSTEM)
			return
		}
		audio := c.state.Current.GetAudioUrl()
		c.players.Play(audio, 0)
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
		c.state.Current.Progress.Progress = c.players.Progress().Progress
		c.Notify(UPDATE_DURATION, SYSTEM)
	}
}

func (c *Controller) AddListner(id string, listener Listener) {
	c.notifications.listners[id] = listener
}

func (c *Controller) RemoveListner(id string) {
	delete(c.notifications.listners, id)
}

func (c *Controller) Notify(action ActionType, user string) {
	state := c.State()
	c.notifications.events <- Event{
		ID: state.ID,
		Action: Action{
			Type: action,
			User: user,
		},
		State: state,
	}
}
