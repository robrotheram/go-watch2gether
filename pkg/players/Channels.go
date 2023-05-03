package players

import (
	"encoding/json"
	"fmt"

	"github.com/asdine/storm"
)

type Store struct {
	Channels map[string]Controller
	*storm.DB
}

func NewStore(path string) (*Store, error) {
	db, err := storm.Open(path)
	store := Store{
		DB:       db,
		Channels: make(map[string]Controller),
	}
	return &store, err
}

func (store *Store) FindChannelById(id string) (*Player, error) {
	var channel Player
	err := store.One("Id", id, &channel)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	if player, ok := store.Channels[id]; ok {
		player.GetState().Active = true
		return player.GetState(), nil
	}
	channel.Active = false
	return &channel, nil
}

func (store *Store) FindControllerById(id string) (Controller, error) {
	var channel Player
	err := store.One("Id", id, &channel)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	if player, ok := store.Channels[id]; ok {
		return player, nil
	}
	return nil, fmt.Errorf("no controller is active")
}

func (store *Store) FindAllChannels() []*Player {
	// players := []*Player{}
	// for _, c := range store.Channels {
	// 	players = append(players, c.GetState())
	// }
	players := []*Player{}
	store.All(&players)
	for _, p := range players {
		p.Active = false
		if _, ok := store.Channels[p.Id]; ok {
			p.Active = true
		}
	}
	return players
}

func (store *Store) RegisterNewChannel(id string, player Controller) error {
	player.SetStore(store.DB)
	if _, found := store.Channels[id]; !found {
		store.Channels[id] = player
	}
	return nil
}

func (store *Store) LeaveChannel(id string) error {
	r, err := store.GetChannel(id)
	if err != nil {
		return fmt.Errorf("room %s not active", id)
	}
	r.Done()
	delete(store.Channels, id)
	return nil
}

func (store *Store) GetChannel(id string) (Controller, error) {
	if channel, found := store.Channels[id]; found {
		return channel, nil
	}
	return nil, fmt.Errorf("channel with ID: %s does not exits", id)
}

func (store *Store) GetAll() {
	var channels []Player
	store.All(&channels)
	data, err := json.Marshal(&channels)
	fmt.Println(err)
	fmt.Println(string(data))
}
