package controllers

import (
	"math/rand"
	"w2g/pkg/media"
)

type PlayState string

var (
	PLAY  = PlayState("PLAY")
	PAUSE = PlayState("PAUSED")
	STOP  = PlayState("STOPPED")
)

type PlayerState struct {
	ID      string        `json:"id"`
	State   PlayState     `json:"status"`
	Queue   []media.Media `json:"queue"`
	Current media.Media   `json:"current"`
	Loop    bool          `json:"loop"`
	Active  bool          `json:"active"`
}

func (state *PlayerState) Next() {
	if len(state.Queue) > 0 {
		state.Current = state.Queue[0]
		state.Queue = state.Queue[1:]
	} else {
		state.Current = media.Media{}
	}
	state.Current.Refresh()
}

func (state *PlayerState) Shuffle() {
	rand.Shuffle(len(state.Queue), func(i, j int) {
		state.Queue[i], state.Queue[j] = state.Queue[j], state.Queue[i]
	})
}

func (state *PlayerState) Repeat() {
	state.Loop = !state.Loop
}

func (state *PlayerState) Reorder(pos1 int, pos2 int) {
	state.Queue = move(state.Queue, pos1, pos2)
}

func (state *PlayerState) Remove(pos1 int) {
	state.Queue = remove(state.Queue, pos1)
}

func (state *PlayerState) Add(videos []media.Media) {
	state.Queue = append(state.Queue, videos...)
}

func (state *PlayerState) Clear() {
	state.Queue = []media.Media{}
}

func (state *PlayerState) ChangeState(ps PlayState) {
	state.State = ps
}

func insert(array []media.Media, value media.Media, index int) []media.Media {
	return append(array[:index], append([]media.Media{value}, array[index:]...)...)
}

func remove(array []media.Media, index int) []media.Media {
	return append(array[:index], array[index+1:]...)
}

func move(array []media.Media, srcIndex int, dstIndex int) []media.Media {
	value := array[srcIndex]
	return insert(remove(array, srcIndex), value, dstIndex)
}
