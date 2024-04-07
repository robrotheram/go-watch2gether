package controllers

import (
	"sync"
	"time"
	"w2g/pkg/media"

	log "github.com/sirupsen/logrus"
)

type PlayerType string
type PlayerMeta struct {
	Id       string              `json:"id"`
	User     string              `json:"user"`
	Progress media.MediaDuration `json:"progress"`
	Type     PlayerType          `json:"type"`
	Running  bool                `json:"running"`
}
type PlayerExitCode uint8

const (
	STOP_EXITCODE PlayerExitCode = iota
	EXIT_EXITCODE
	SKIP_EXITCODE
)

type Player interface {
	Play(string, int) (PlayerExitCode, error)
	Seek(time.Duration)
	Pause()
	Unpause()
	Stop()
	Close()
	Meta() PlayerMeta
	Id() string
}

type Players struct {
	players  map[string]Player
	AutoSkip bool
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
	p.players[player.Meta().Id] = player
}

func (p *Players) Remvoe(id string) {
	if player, ok := p.players[id]; ok {
		player.Close()
		delete(p.players, id)
	}
}

func (p *Players) Seek(seconds time.Duration) {
	for _, player := range p.players {
		player.Seek(seconds)
	}
}

func (p *Players) GetProgress() []PlayerMeta {
	data := []PlayerMeta{}
	for _, player := range p.players {
		data = append(data, player.Meta())
	}
	return data
}

func (p *Players) Progress() media.MediaDuration {
	progress := media.MediaDuration{}
	for _, player := range p.players {
		prg := player.Meta().Progress.Progress
		if prg > progress.Progress {
			progress.Progress = prg
		}
	}
	return progress
}

func (p *Players) Play(url string, start int) {
	wg := sync.WaitGroup{}
	wg.Add(len(p.players))
	for _, player := range p.players {
		player := player //TODO: Remove in when we upgrade go 1.22
		go func(player Player) {
			exit, err := player.Play(url, start)
			if err != nil {
				log.Warnf("%s player error: %v", player.Meta().Type, err)
			}
			wg.Done()
			if exit == STOP_EXITCODE {
				p.Stop()
			}
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
		if player.Meta().Running {
			return true
		}
	}
	return false
}
