package controllers

import (
	"context"
	"sync"
	"time"
	"w2g/pkg/history"
	"w2g/pkg/media"
	"w2g/pkg/utils"

	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
	bolt "go.etcd.io/bbolt"
)

const SYSTEM = "system"

type Controller struct {
	state         *PlayerState
	players       *Players
	notifications *Notify
	running       bool
	db            *utils.Store[*PlayerState]
	history       *history.HistoryStore
	sync          *sync.Mutex
}

func NewController(id string, db *bolt.DB) *Controller {
	store := &utils.Store[*PlayerState]{
		DB:     db,
		Bucket: []byte("controllers"),
	}
	store.Create()
	contoller := Controller{
		players:       newPlayers(),
		notifications: NewNotifications(),
		history:       history.NewHistoryStore(db),
		db:            store,
		sync:          &sync.Mutex{},
	}
	contoller.load(id)
	contoller.AddListner(uuid.NewString(), &Auditing{})
	if utils.Configuration.BetterStackToken != "" {
		log.Info("enabling better stack logging")
		contoller.AddListner(uuid.NewString(), &BetterStack{Token: utils.Configuration.BetterStackToken})
	}

	return &contoller
}

func (c *Controller) save() error {
	return c.db.Save(c.state.ID, c.state)
}

func (c *Controller) load(id string) {
	state, err := c.db.Get(id)
	if err != nil {
		state = &PlayerState{
			ID:    id,
			Queue: []*media.Media{},
			State: STOP,
			Loop:  false,
		}
		c.db.Save(id, state)
	}
	state.ChangeState(STOP)
	c.state = state
}

func (c *Controller) Start(user string) {
	if c.players.Empty() {
		return
	}
	if !c.running {
		if c.state.Current == nil {
			c.state.Next()
		}
		c.running = true
		go c.progress()
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
	c.state.ChangeState(STOP)
	c.players.Stop()
	c.save()
	c.Notify(STOP_ACTION, user)
}

func (c *Controller) Pause(user string) {
	c.players.Pause()
	c.state.ChangeState(PAUSE)
	c.save()
	c.Notify(PAUSE_ACTION, user)
}

func (c *Controller) Seek(seconds time.Duration, user string) {
	c.players.Seek(seconds)
	c.state.Current.Progress.Progress = c.players.Progress().Progress
	c.Notify(SEEK, user)
}

func (c *Controller) Add(url string, top bool, user string) error {
	tracks, err := media.NewVideo(url, user)
	go media.RefreshAll(tracks)
	if err != nil {
		return err
	}
	if top {
		c.state.AddTop(tracks)
	} else {
		c.state.AddBottom(tracks)
	}
	c.save()
	c.Notify(ADD_QUEUE, user)
	return nil
}

func (c *Controller) Skip(user string) {
	if c.running {
		c.players.Stop()
	} else {
		c.state.Next()
	}
	c.Notify(SKIP_ACTION, user)
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

func (c *Controller) UpdateQueue(videos []*media.Media, user string) {
	c.state.Queue = videos
	c.save()
	c.Notify(UPDATE_QUEUE, user)
}

func (c *Controller) State() *PlayerState {
	c.state.Active = !c.players.Empty()
	return c.state
}

func (c *Controller) History() ([]*media.Media, error) {
	return c.history.GetHisory(c.state.ID)
}

func (c *Controller) ServerState() ServerState {
	return ServerState{
		Players: c.players.GetProgress(),
		State:   *c.state,
	}
}

func (c *Controller) Update(state *PlayerState, user string) {
	c.state = state
	c.save()
	c.Notify(UPDATE, user)
}

func (c *Controller) Join(player Player, user string) {
	if _, ok := c.players.players[player.Id()]; !ok {
		c.players.Add(player)
		c.Notify(PLAYER_ACTION, user)
	}
}

func (c *Controller) Leave(id string, user string) {
	c.players.Remvoe(id)
	c.Notify(LEAVE_ACTION, user)
	if len(c.players.players) == 0 {
		c.Stop(SYSTEM)
	}
}

func (c *Controller) ContainsPlayer(id string) bool {
	if _, ok := c.players.players[id]; ok {
		return true
	}
	return false
}

func (c *Controller) progress() {
	for {
		if c.state.Current == nil {
			return
		}
		audio := c.state.Current.AudioUrl
		c.history.AddTrack(c.state.ID, c.state.Current)
		log.Debug("START_PLAYING")
		ctx, cancel := context.WithCancel(context.Background())
		go c.duration(ctx)
		c.players.Play(audio, 0)
		log.Debug("STOP_PLAYING")
		cancel()
		if !c.state.Loop {
			c.state.Next()
			c.Notify(UPDATE_QUEUE, SYSTEM)
		}
		if c.state.Current == nil || c.players.Empty() || c.state.State == STOP {
			c.Stop(SYSTEM)
			log.Debug("DONE")
			return
		}
		log.Debug("NEXT")
	}
}

func (c *Controller) duration(ctx context.Context) {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop() // Ensure ticker is cleaned up when the goroutine exits
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			if c.state.Current != nil {
				c.state.Current.Progress.Progress = c.players.Progress().Progress
				c.Notify(UPDATE_DURATION, SYSTEM)
			}
		}
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
			Type:    action,
			User:    user,
			Channel: state.ID,
		},
		State:   *state,
		Players: c.players.GetProgress(),
	}
}

func (c *Controller) Players() *Players {
	return c.players
}
