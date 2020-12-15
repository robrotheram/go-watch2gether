package pkg

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/segmentio/ksuid"
)

type Video struct {
	UID  string `json:"uid"`
	User string `json:"user"`
	Url  string `json:"url"`
}

type Event struct {
	Action       string  `json:"action"`
	Host         string  `json:"host"`
	User         User    `json:"user"`
	Queue        []Video `json:"queue"`
	CurrentVideo Video   `json:"current_video"`
	Seek         float32 `json:"seek"`
	Users        []User  `json:"users"`
	Controls     bool    `json:"controls"`
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
	Name         string `json:"name"`
	Host         string `json:"host"`
	PreVideo     Video
	CurrentVideo Video   `json:"current_video"`
	Seek         float32 `json:"seek"`
	Controls     bool    `json:"controls"`
	Playing      bool    `json:"playing"`
	Queue        []Video `json:"queue"`
	Users        []User  `json:"users"`
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
	//fmt.Println(data)
	switch data.Action {

	case "PLAYING":
		r.Meta.Playing = true
		break
	case "PAUSING":
		r.Meta.Playing = false
		break
	case "UPDATE_HOST":
		r.updateHost(User{Name: data.Host})
		break
	case "NEXT_VIDEO":
		r.skipNextVideo()
		return false
	case "SEEK":
		r.Meta.Seek = data.Seek
		break
	case "UPDATE_CONTROLS":
		r.Meta.Controls = data.Controls
	case "SEEK_TO_ME":
		r.SeekToUser(data.User)
		return false
	case "UPDATE_QUEUE":
		return r.handleQueueUpdate(data.Queue)
	case "HANDLE_FINSH":
		r.HandleFinish(data.User)
		return false
	case "USER_UPADTE":
		r.SeenUser(data.User)
		//Do not sink with client
		return false
	}
	return true
}

func (room *room) handleQueueUpdate(queue []Video) bool {
	room.Meta.Queue = queue
	if room.Meta.CurrentVideo.UID == "" {
		room.skipNextVideo()
		return false
	}
	return true
}

func (room *room) updateHost(user User) {
	room.Meta.Host = user.Name
	for i := range room.Meta.Users {
		if room.Meta.Users[i].Name == user.Name {
			room.Meta.Users[i].IsHost = true
		}
	}
}

func (r *room) skipNextVideo() {
	if len(r.Meta.Queue) == 0 {
		r.Meta.PreVideo = r.Meta.CurrentVideo
		r.Meta.CurrentVideo = Video{}
	} else {
		video := r.Meta.Queue[0]
		r.Meta.Queue = r.Meta.Queue[1:]
		r.Meta.PreVideo = r.Meta.CurrentVideo
		r.Meta.CurrentVideo = video
	}

	evt := Event{
		Action:       "CHANGE_VIDEO",
		CurrentVideo: r.Meta.CurrentVideo,
		Queue:        r.Meta.Queue,
	}

	b, _ := json.Marshal(evt)
	r.forward <- b
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
		Meta:    roomMeta{Name: name, Users: []User{}, Queue: []Video{}, Controls: true},
		status:  "Initilised",
		ID:      ksuid.New().String(),
	}
}

func (room *room) Join(user User) {
	fmt.Println(user)
	if len(room.Meta.Users) == 0 {
		room.Meta.Host = user.Name
		user.IsHost = true
	}
	room.Meta.Users = append(room.Meta.Users, user)

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
		if user.LastSeen.Add(10 * time.Second).Before(time.Now()) {
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
		room.updateHost(room.Meta.Users[0])
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

func (room *room) UpdateUserSeek(u User) {
	for i := range room.Meta.Users {
		user := &room.Meta.Users[i]
		if u.Name == user.Name {
			user.Seek = u.Seek
		}
	}
	if room.Meta.Host == u.Name {
		room.Meta.Seek = u.Seek
	}
}
func (room *room) SeekToUser(u User) {
	for i := range room.Meta.Users {
		user := &room.Meta.Users[i]
		if u.Name == user.Name {
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

func (room *room) SeenUser(u User) {
	found := false
	for i := range room.Meta.Users {
		user := &room.Meta.Users[i]

		if u.Name == user.Name {
			if user.IsHost {
				room.Meta.Seek = u.Seek
			}
			user.LastSeen = time.Now()
			user.Seek = u.Seek
			user.CurrentVideo = u.CurrentVideo
			found = true
			room.SendUsersProgress()
			return
		}

	}
	if !found {
		room.Join(u)
	}
}

func (room *room) HandleFinish(user User) {

	if room.Meta.PreVideo.UID == user.CurrentVideo.UID {
		return
	}
	for i := range room.Meta.Users {
		u := &room.Meta.Users[i]
		if u.Seek < float32(1) && u.Name != user.Name {
			return
		} else if u.Name == user.Name {
			u.Seek = float32(1)
		}
	}
	room.skipNextVideo()
}

func (room *room) SendUsersProgress() {
	evt := Event{
		Action: "ON_PROGRESS_UPDATE",
		Users:  room.Meta.Users,
	}
	b, _ := json.Marshal(evt)
	room.forward <- b
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
	if m.CurrentVideo != meta.CurrentVideo && meta.CurrentVideo.Url != "" {
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
