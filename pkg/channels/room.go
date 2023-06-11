package channels

type RoomType int
type EventType int

const (
	DISCORD RoomType = iota
	WEBSOCKET
)

const (
	PLAY EventType = iota
	PAUSE
	STOP
	SKIP
)

type Room struct {
	controller map[RoomType]Controller
	noify      chan EventType
	quit       chan bool
}

func (room *Room) NewChannel(cntr Controller) {
	room.controller[cntr.GetType()] = cntr
}

func NewRoom(cntr Controller) *Room {
	room := &Room{
		noify:      make(chan EventType),
		quit:       make(chan bool),
		controller: map[RoomType]Controller{},
	}
	room.NewChannel(cntr)
	go func() {
		for {
			select {
			case <-room.quit:
				return
			case message := <-room.noify:
				for _, c := range room.controller {
					c.Notify(message)
				}
			}
		}
	}()
	return room
}
