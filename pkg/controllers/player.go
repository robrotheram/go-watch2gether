package controllers

import (
	"sync"
	"w2g/pkg/media"
)

type PlayerType string
type Player interface {
	Play(*sync.WaitGroup, string, int) error
	Progress() media.MediaDuration
	Type() PlayerType
	Pause()
	Unpause()
	Stop()
	Close()
}

type Players struct {
	players map[PlayerType]Player
}

func newPlayers() *Players {
	return &Players{
		players: map[PlayerType]Player{},
	}
}

func (p *Players) Empty() bool {
	return len(p.players) == 0
}

func (p *Players) Add(player Player) {
	p.players[player.Type()] = player
}

func (p *Players) Remvoe(id PlayerType) {
	p.players[id].Close()
	// delete(p.players, id)
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
		go player.Play(&wg, url, start)
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
