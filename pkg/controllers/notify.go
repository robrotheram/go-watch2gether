package controllers

type ActionType string

type Action struct {
	ActionType ActionType `json:"type"`
	User       string     `json:"user"`
}

var (
	PLAY_ACTION     = ActionType("PLAY")
	PAUSE_ACTION    = ActionType("PAUSE")
	ADD_QUEUE       = ActionType("ADD_QUEUE")
	UPDATE_QUEUE    = ActionType("UPDATE_QUEUE")
	UPDATE          = ActionType("UPDATE")
	REMOVE_QUEUE    = ActionType("REMOVE_QUEUE")
	UPDATE_DURATION = ActionType("UPDATE_DURATION")
	STOP_ACTION     = ActionType("STOP")
	LOOP_ACTION     = ActionType("LOOP")
	SHUFFLE_ACTION  = ActionType("SHUFFLE")
	SKIP_ACTION     = ActionType("SKIP")
	PLAYER_ACTION   = ActionType("PlAYER CHANGE")
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
