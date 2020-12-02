package pkg

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/segmentio/ksuid"
)

type VideoQueue struct {
	User string `json:"user"`
	Url  string `json:"url"`
}

type Event struct {
	Action       string       `json:"action"`
	Host         string       `json:"host"`
	User         string       `json:"user"`
	Queue        []VideoQueue `json:"queue"`
	CurrentVideo string       `json:"current_video"`
	Seek         float32      `json:"seek"`
	Users        []User       `json:"users"`
}

type room struct {
	// forward is a channel that holds incoming messages
	// that should be forwarded to the other clients.
	forward chan []byte
	// join is a channel for clients wishing to join the room.
	join chan *client
	// leave is a channel for clients wishing to leave the room.
	leave chan *client

	//Channel to quit the room
	quit   chan bool
	status string
	// clients holds all current clients in this room.
	clients map[*client]bool
	ID      string
	Meta    roomMeta
}

type roomMeta struct {
	Name         string       `json:"name"`
	Host         string       `json:"host"`
	CurrentVideo string       `json:"current_video"`
	Seek         float32      `json:"seek"`
	Controls     bool         `json:"controls"`
	Playing      bool         `json:"playing"`
	Queue        []VideoQueue `json:"queue"`
	Users        []User       `json:"users"`
}

func (r *room) run() {
	r.status = "Running"
	for {
		select {
		case <-r.quit:
			// kill the goroutine
			r.status = "Stopping"
			for client := range r.clients {
				delete(r.clients, client)
			}
			r.status = "Stopped"
			return
		case client := <-r.join:
			// joining
			r.clients[client] = true
		case client := <-r.leave:
			// leaving
			delete(r.clients, client)
			close(client.send)
		case msg := <-r.forward:
			// forward message to all clients
			if r.processEvent(msg) {
				for client := range r.clients {
					client.send <- msg
				}
			}

		}
	}
}
func (r *room) processEvent(byteData []byte) bool {
	var data Event
	in := bytes.NewReader(byteData)
	_ = json.NewDecoder(in).Decode(&data)
	fmt.Println(data)
	switch data.Action {

	case "PLAYING":
		r.Meta.Playing = true
		break
	case "PAUSING":
		r.Meta.Playing = false
		break
	case "UPDATE_HOST":
		r.Meta.Host = data.Host
		break
	case "CHANGE_VIDEO":
		r.Meta.CurrentVideo = data.CurrentVideo
		break
	case "SEEK":
		r.Meta.Seek = data.Seek
		break
	case "SEEK_TO_ME":
		r.SeekToUser(data.User)
		return false
	case "UPDATE_QUEUE":
		r.Meta.Queue = data.Queue
		break
	case "ON_PROGRESS_UPDATE":
		r.UpdateUserSeek(data.User, data.Seek)
		//Do not sink with client
		return false
	case "USER_UPATE":
		r.SeenUser(data.User)
		//Do not sink with client
		return false
	}
	return true
}

func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
}

func (r *room) ServeHTTP(w http.ResponseWriter, req *http.Request) error {
	enableCors(&w)
	socket, err := upgrader.Upgrade(w, req, nil)
	if err != nil {
		return StatusError{http.StatusBadRequest, fmt.Errorf("Connection is not using the websocket protocol")}
	}
	client := &client{
		socket: socket,
		send:   make(chan []byte, messageBufferSize),
		room:   r,
	}
	r.join <- client
	defer func() { r.leave <- client }()
	go client.write()
	client.read()
	return nil
}

func newRoom(name string) *room {
	return &room{
		forward: make(chan []byte),
		join:    make(chan *client),
		leave:   make(chan *client),
		clients: make(map[*client]bool),
		quit:    make(chan bool),
		Meta:    roomMeta{Name: name, Users: []User{}, Queue: []VideoQueue{}, Controls: true},
		status:  "Initilised",
		ID:      ksuid.New().String(),
	}
}

func (room *room) Join(username string) {
	if len(room.Meta.Users) == 0 {
		room.Meta.Host = username
	}
	room.Meta.Users = append(room.Meta.Users, NewUser(username))

	evt := Event{
		Action: "USER_UPDATED",
		Users:  room.Meta.Users,
	}
	b, _ := json.Marshal(evt)
	room.forward <- b
}

func (room *room) PurgeUsers() {
	deletedUsers := []string{}

	for i := range room.Meta.Users {
		user := &room.Meta.Users[i]
		if user.LastSeen.Add(5 * time.Second).Before(time.Now()) {
			deletedUsers = append(deletedUsers, user.Name)
		}
	}
	for _, user := range deletedUsers {
		room.Leave(user)
	}
}

func (room *room) Leave(username string) {
	for i, v := range room.Meta.Users {
		if v.Name == username {
			room.Meta.Users = append(room.Meta.Users[:i], room.Meta.Users[i+1:]...)
			break
		}
	}

	if room.Meta.Host == username {
		if len(room.Meta.Users) == 0 {
			return
		}
		room.Meta.Host = room.Meta.Users[0].Name
		evt := Event{
			Action: "UPDATE_HOST",
			Host:   room.Meta.Host,
		}
		b, _ := json.Marshal(evt)
		room.forward <- b

	}

	evt := Event{
		Action: "USER_UPDATED",
		Users:  room.Meta.Users,
	}
	b, _ := json.Marshal(evt)
	room.forward <- b
}

func (room *room) UpdateUserSeek(username string, seek float32) {
	for i := range room.Meta.Users {
		user := &room.Meta.Users[i]
		if username == user.Name {
			user.Seek = seek
		}
	}
	if room.Meta.Host == username {
		room.Meta.Seek = seek
	}
}
func (room *room) SeekToUser(username string) {
	for i := range room.Meta.Users {
		user := &room.Meta.Users[i]
		if username == user.Name {
			room.Meta.Seek = user.Seek
		}
	}

	evt := Event{
		Action: "SEEK_TO_USER",
		Seek:   room.Meta.Seek,
	}
	b, _ := json.Marshal(evt)
	room.forward <- b
}
func (room *room) SeenUser(username string) {
	for i := range room.Meta.Users {
		user := &room.Meta.Users[i]
		if username == user.Name {
			user.LastSeen = time.Now()
		}
	}
}

func (room *room) ContainsUser(name string) bool {
	for _, user := range room.Meta.Users {
		if name == user.Name {
			return true
		}
	}
	return false
}

func (m *roomMeta) Update(meta roomMeta) {

	if m.Name != meta.Name && meta.Name != "" {
		m.Name = meta.Name
	}
	if m.Host != meta.Host && meta.Host != "" {
		m.Host = meta.Host
	}
	if m.CurrentVideo != meta.CurrentVideo && meta.CurrentVideo != "" {
		m.CurrentVideo = meta.CurrentVideo
	}
	if m.Seek != meta.Seek {
		m.Seek = meta.Seek
	}
	if m.Controls != meta.Controls {
		m.Controls = meta.Controls
	}
	if m.Playing != meta.Playing {
		m.Playing = meta.Playing
	}
	if meta.Queue != nil {
		m.Queue = meta.Queue
	}
}
