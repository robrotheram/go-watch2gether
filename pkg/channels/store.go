package channels

import (
	"fmt"
	"watch2gether/pkg/channels/model"
	"watch2gether/pkg/channels/players"

	"github.com/asdine/storm"
)

type Store struct {
	Channels map[string]*Controller
	*storm.DB
}

func NewStore(path string) (*Store, error) {
	db, err := storm.Open(path)
	store := Store{
		DB:       db,
		Channels: make(map[string]*Controller),
	}
	return &store, err
}

func (store *Store) GetController(id string) (*Controller, error) {
	if channel, found := store.Channels[id]; found {
		return channel, nil
	}
	return nil, fmt.Errorf("channel with ID: %s does not exits", id)
}

func (store *Store) GetState(id string) (*model.PlayerState, error) {
	ctrl := Controller{
		Id: id,
		DB: store.DB,
	}
	return ctrl.GetState()
}

func (store *Store) Register(id string, pType players.PlayerType, player players.Player) error {
	if channel, found := store.Channels[id]; found {
		channel.WithPlayer(pType, player)
		return nil
	}
	store.Channels[id] = NewContoller(id, store.DB).WithPlayer(pType, player)

	return nil
}

func (store *Store) RemovePlayer(id string, pType players.PlayerType) error {
	if channel, found := store.Channels[id]; found {
		channel.RemovePlayer(pType)
		return nil
	}
	return fmt.Errorf("channel with ID: %s does not exits", id)
}

func (store *Store) FindAllChannels() []*model.PlayerState {
	players := []*model.PlayerState{}
	store.All(&players)
	for _, p := range players {
		p.Active = false
		if _, ok := store.Channels[p.Id]; ok {
			p.Active = true
		}
	}
	return players
}
