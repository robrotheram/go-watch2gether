package controllers

import (
	"math/rand"
	"w2g/pkg/media"
)

type PlayerState struct {
	Id      string
	State   string
	Queue   []media.Media
	Current media.Media
}

func (state *PlayerState) Next() {
	if len(state.Queue) > 0 {
		state.Current = state.Queue[0]
		state.Queue = state.Queue[1:]
	} else {
		state.Current = media.Media{}
	}
}

func (state *PlayerState) Shuffle() {
	rand.Shuffle(len(state.Queue), func(i, j int) {
		state.Queue[i], state.Queue[j] = state.Queue[j], state.Queue[i]
	})
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
