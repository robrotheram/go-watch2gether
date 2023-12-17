package controllers

type Action string

var (
	PLAY_ACTION        = Action("PLAY")
	PAUSE_ACTION       = Action("PAUSE")
	UPDATEQUEUE_ACTION = Action("UPDATE_QUEUE")
	STOP_ACTION        = Action("STOP")
	LOOP_ACTION        = Action("LOOP")
	SHUFFLE_ACTION     = Action("SHUFFLE")
	SKIP_ACTION        = Action("SKIP")
	PLAYER_ACTION      = Action("PlAYER CHANGE")
)

type Event struct {
	ID     string
	Action Action
	State  PlayerState
}

type Listener interface {
	Send(Event)
}

type Notify struct {
	events   chan Event
	done     chan any
	listners []Listener
}

func NewNotifications() *Notify {
	notify := Notify{
		events:   make(chan Event),
		done:     make(chan any),
		listners: []Listener{},
	}
	go notify.start()
	return &notify
}

func (n *Notify) Close() {
	n.done <- "close"
}

func (n *Notify) start() {
	for {
		select {
		case <-n.done:
			return
		case event := <-n.events:
			for _, listener := range n.listners {
				listener.Send(event)
			}
		}
	}
}
