package controllers

type Hub struct {
	channels map[string]*Controller
}

func NewHub() *Hub {
	return &Hub{
		channels: make(map[string]*Controller),
	}
}

func (hub *Hub) Add(id string, controller *Controller) {
	hub.channels[id] = controller
}

func (hub *Hub) Get(id string) *Controller {
	if _, ok := hub.channels[id]; !ok {
		hub.channels[id] = NewController()
	}
	return hub.channels[id]
}
