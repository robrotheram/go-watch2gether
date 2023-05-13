package channels

import (
	"time"
	"watch2gether/pkg/media"

	"github.com/asdine/storm"
)

type Controller interface {
	Play()
	Pause()
	Stop()
	Skip()
	Shuffle()
	Done()
	Move(int, int)
	Remove(int)
	Clear()
	SetLoop(bool)
	Duration() time.Duration
	Add(meida []media.Media)
	GetQueue() []media.Media
	UpdateQueue([]media.Media)
	GetCurrentVideo() media.Media
	SetStore(*storm.DB)
	GetState() *Player
}
