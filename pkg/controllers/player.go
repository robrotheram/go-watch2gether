package controllers

import (
	"fmt"
	"sync"
	"time"
	"w2g/pkg/media"
)

type PlayerType string
type PlayerMeta struct {
	Progress media.MediaDuration `json:"progress"`
	Type     PlayerType          `json:"type"`
	Running  bool                `json:"running"`
}
type Player interface {
	Play(string, int) error
	Progress() media.MediaDuration
	Seek(time.Duration)
	Type() PlayerType
	Id() string
	Pause()
	Unpause()
	Stop()
	Close()
	Status() bool
}

type Players struct {
	players map[string]Player
}

func newPlayers() *Players {
	return &Players{
		players: map[string]Player{},
	}
}

func (p *Players) Empty() bool {
	return len(p.players) == 0
}

func (p *Players) Add(player Player) {
	p.players[player.Id()] = player
}

func (p *Players) Seek(seconds time.Duration) {
	for _, player := range p.players {
		player.Seek(seconds)
	}
}

func (p *Players) GetProgress() map[string]PlayerMeta {
	data := map[string]PlayerMeta{}
	for _, player := range p.players {
		data[player.Id()] = PlayerMeta{
			Progress: player.Progress(),
			Type:     player.Type(),
			Running:  player.Status(),
		}
	}
	return data
}

func (p *Players) Remvoe(id string) {
	if player, ok := p.players[id]; ok {
		player.Close()
		delete(p.players, id)
	}
}

func (p *Players) Progress() media.MediaDuration {
	progress := media.MediaDuration{}
	for _, player := range p.players {
		prg := player.Progress().Progress
		if prg > progress.Progress {
			progress.Progress = prg
		}
	}
	return progress
}

func (p *Players) Play(url string, start int) {
	wg := sync.WaitGroup{}
	for _, player := range p.players {
		wg.Add(1)
		go func(player Player) {
			err := player.Play(url, start)
			if err != nil {
				fmt.Printf("%s player error: %v", player.Type(), err)
			}
			wg.Done()
			//Currently Set it the the first "player to finish will override all other players"
			p.Stop()
		}(player)
	}
	wg.Wait()
}

func (p *Players) Pause() {
	for _, player := range p.players {
		player.Pause()
	}
}

func (p *Players) Unpause() {
	for _, player := range p.players {
		player.Unpause()
	}
}

func (p *Players) Stop() {
	for _, player := range p.players {
		player.Stop()
	}
}

func (p *Players) Close() {
	for _, player := range p.players {
		player.Close()
	}
}

func (p *Players) Running() bool {
	for _, player := range p.players {
		if player.Status() {
			return true
		}
	}
	return false
}
