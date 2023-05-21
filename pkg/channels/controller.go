package channels

import (
	"math/rand"
	"sync"
	"time"
	"watch2gether/pkg/media"

	"github.com/asdine/storm"
)

type Controller interface {
	Play() error
	Pause() error
	Stop() error
	Skip() error
	Shuffle() error
	Run() error
	Done() error
	Move(int, int) error
	Remove(int) error
	Clear() error
	SetLoop(bool) error
	Duration() time.Duration
	Add(meida []media.Media) error
	UpdateQueue([]media.Media) error
	GetQueue() ([]media.Media, error)
	GetCurrentVideo() (media.Media, error)
	SetStore(*storm.DB)
	GetState() (*Player, error)
}

type Base struct {
	id string
	*storm.DB
	*sync.Mutex
}

func NewBase(id string) *Base {
	return &Base{
		id:    id,
		Mutex: &sync.Mutex{},
	}
}

func (base *Base) update(state *Player) error {
	return base.Save(state)
}

func (base *Base) SetStore(db *storm.DB) {
	base.DB = db
	if _, err := base.GetState(); err != nil {
		db.Save(&Player{
			Id:    base.id,
			State: STOPPED,
			Queue: []media.Media{},
		})
	}
}

func (base *Base) GetState() (*Player, error) {
	var state Player
	err := base.One("Id", base.id, &state)
	return &state, err
}

func (base *Base) UpdateQueue(videos []media.Media) error {
	base.Lock()
	defer base.Unlock()

	state, err := base.GetState()
	if err != nil {
		return err
	}
	state.Queue = videos
	return base.update(state)
}

func (base *Base) SetLoop(loop bool) error {
	base.Lock()
	defer base.Unlock()
	state, err := base.GetState()
	if err != nil {
		return err
	}
	state.Loop = loop
	return base.update(state)
}

func (base *Base) Add(meida []media.Media) error {
	base.Lock()
	defer base.Unlock()
	state, err := base.GetState()
	if err != nil {
		return err
	}
	state.Queue = append(state.Queue, meida...)
	return base.update(state)
}

func (base *Base) Next() error {
	base.Lock()
	defer base.Unlock()

	state, err := base.GetState()
	if err != nil {
		return err
	}
	state.Current = media.Media{}
	if len(state.Queue) > 0 {
		state.Current, state.Queue = state.Queue[0], state.Queue[1:]
	}
	return base.update(state)
}

func (base *Base) Shuffle() error {
	base.Lock()
	defer base.Unlock()

	state, err := base.GetState()
	if err != nil {
		return err
	}
	rand.Shuffle(len(state.Queue), func(i, j int) {
		state.Queue[i], state.Queue[j] = state.Queue[j], state.Queue[i]
	})
	return base.update(state)
}

func (base *Base) GetQueue() ([]media.Media, error) {
	base.Lock()
	defer base.Unlock()
	state, err := base.GetState()
	if err != nil {
		return []media.Media{}, err
	}
	return state.Queue, nil
}

func (base *Base) GetCurrentVideo() (media.Media, error) {
	state, err := base.GetState()
	if err != nil {
		return media.Media{}, err
	}
	return state.Current, nil
}

func (base *Base) Clear() error {
	base.Lock()
	defer base.Unlock()

	state, err := base.GetState()
	if err != nil {
		return err
	}
	state.Queue = []media.Media{}
	return base.update(state)
}

func (base *Base) Move(srcIndex int, dstIndex int) error {
	state, err := base.GetState()
	if err != nil {
		return err
	}
	state.Move(srcIndex, dstIndex)
	return base.update(state)
}

func (base *Base) Remove(srcIndex int) error {
	state, err := base.GetState()
	if err != nil {
		return err
	}
	state.Remove(srcIndex)
	return base.update(state)
}
