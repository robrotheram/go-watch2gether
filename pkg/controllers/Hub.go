package controllers

import "fmt"

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

func (hub *Hub) New(id string) *Controller {
	hub.channels[id] = NewController()
	return hub.channels[id]
}

func (hub *Hub) Get(id string) (*Controller, error) {
	if _, ok := hub.channels[id]; !ok {
		return nil, fmt.Errorf("channel not found")
	}
	return hub.channels[id], nil
}
