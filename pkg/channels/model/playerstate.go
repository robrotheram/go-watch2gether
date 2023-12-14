package model

import (
	"watch2gether/pkg/media"
)

type State string

var (
	PLAYING = State("PLAYING")
	PAUSED  = State("PAUSED")
	STOPPED = State("STOPPED")
)

type PlayerState struct {
	Id        string `storm:"id" json:"id"`
	Loop      bool
	Active    bool
	IsPlaying bool
	Queue     []media.Media
	Current   media.Media
}

func NewPlayerState(id string) *PlayerState {
	return &PlayerState{
		Id:        id,
		IsPlaying: false,
		Queue:     []media.Media{},
	}
}

func (p *PlayerState) Next() {
	p.Current = media.Media{}
	if len(p.Queue) > 0 {
		p.Current, p.Queue = p.Queue[0], p.Queue[1:]
	}
}

func insert(array []media.Media, value media.Media, index int) []media.Media {
	return append(array[:index], append([]media.Media{value}, array[index:]...)...)
}

func (p *PlayerState) Remove(index int) []media.Media {
	return append(p.Queue[:index], p.Queue[index+1:]...)
}

func (p *PlayerState) Move(srcIndex int, dstIndex int) []media.Media {
	value := p.Queue[srcIndex]
	return insert(p.Remove(srcIndex), value, dstIndex)
}

func (p *PlayerState) RemoveDuplicates() {
	// map to store unique keys
	keys := make(map[string]bool)
	returnSlice := []media.Media{}
	for _, item := range p.Queue {
		if _, value := keys[item.ID]; !value {
			keys[item.ID] = true
			returnSlice = append(returnSlice, item)
		}
	}
	p.Queue = returnSlice
}
