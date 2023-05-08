package players

import (
	"time"
	"watch2gether/pkg/media"

	"github.com/asdine/storm"
)

var (
	PLAYING = PlayerState("PLAYING")
	PAUSED  = PlayerState("PAUSED")
	STOPPED = PlayerState("STOPPED")
)

type Controller interface {
	Play()
	Pause()
	Stop()
	Skip()
	Shuffle()
	Done()
	Load()
	Clear()
	SetLoop(bool)
	Duration() time.Duration
	Add(meida []media.Media)
	GetQueue() []media.Media
	UpdateQueue([]media.Media)
	GetState() *Player
	UpdaetState(*Player)
	GetCurrentVideo() media.Media
	SetStore(*storm.DB)
}

type PlayerState string

type Player struct {
	Id          string `storm:"id" json:"id"`
	State       PlayerState
	Loop        bool
	Active      bool
	Queue       []media.Media
	Proccessing time.Duration
	Current     media.Media
}

func (p *Player) MediaRefresh() {
	media.RefreshAudioURL(&p.Current)
}

func insert(array []media.Media, value media.Media, index int) []media.Media {
	return append(array[:index], append([]media.Media{value}, array[index:]...)...)
}

func (p *Player) Remove(index int) []media.Media {
	return append(p.Queue[:index], p.Queue[index+1:]...)
}

func (p *Player) Move(srcIndex int, dstIndex int) []media.Media {
	value := p.Queue[srcIndex]
	return insert(p.Remove(srcIndex), value, dstIndex)
}
