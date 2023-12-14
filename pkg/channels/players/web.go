package players

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
	"watch2gether/pkg/channels/model"

	"github.com/gorilla/websocket"
)

var WEB = PlayerType("web")

type WebPlayer struct {
	Clients     map[*Client]bool
	Progress    float64
	ProgresTime time.Duration
	event       model.Event
	notify      chan model.Event
}

func (wp *WebPlayer) Start(notify chan model.Event) {
	wp.notify = notify
}

func NewWebPlayer(id string) *WebPlayer {
	player := WebPlayer{
		Clients: map[*Client]bool{},
	}
	return &player
}

func (wp *WebPlayer) Join(client *Client) {
	wp.Clients[client] = true
}

func (wp *WebPlayer) recive(msg []byte) {
	var evt model.Event
	err := json.Unmarshal(msg, &evt)
	if err != nil {
		return
	}
	fmt.Println(evt.Action)
	switch evt.Action {
	case model.UPADATE:
		// if (float64(evt.State.Proccessing)) > 0 {
		// 	wp.Progress = float64(evt.State.Proccessing) / float64(evt.State.Current.Duration.Nanoseconds())
		// 	wp.ProgresTime = evt.State.Proccessing
		// }
	case model.FINISHED:
		wp.onFinish()
	}

}

func (wp *WebPlayer) onFinish() {
	wp.Progress = 1
	wp.notify <- wp.event.WithAction(model.FINISHED)
}

func (wp *WebPlayer) leave(client *Client) {
	delete(wp.Clients, client)
}

func (wp *WebPlayer) Notify(event model.Event) {
	for client := range wp.Clients {
		data, _ := json.Marshal(event)
		client.send <- data
	}
}

func (wp *WebPlayer) Pause() error {
	return nil
}

func (wp *WebPlayer) Stop() error {
	return nil
}

func (wp *WebPlayer) Skip() error {
	wp.onFinish()
	return nil
}

func (wp *WebPlayer) Quit() error {
	return nil
}

func (wp *WebPlayer) Play(event model.Event) error {
	wp.event = event
	wp.Progress = 0
	wp.ProgresTime = 0
	return nil
}

func (wp *WebPlayer) Duration() time.Duration {
	return wp.ProgresTime
}

type Client struct {
	socket *websocket.Conn
	send   chan []byte
	player *WebPlayer
}

const (
	SocketBufferSize  = 1024
	MessageBufferSize = 256
)

var Upgrader = &websocket.Upgrader{
	ReadBufferSize:  SocketBufferSize,
	WriteBufferSize: SocketBufferSize,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func (c *Client) Read() {
	defer c.socket.Close()
	for {
		_, msg, err := c.socket.ReadMessage()
		if err != nil {
			fmt.Printf("ERROR decoding %v", err)
			return
		}
		c.player.recive(msg)
	}
}

func (c *Client) Write() {
	defer c.socket.Close()
	for msg := range c.send {
		err := c.socket.WriteMessage(websocket.TextMessage, msg)
		if err != nil {
			return
		}
	}
}

func NewClient(r *WebPlayer, socket *websocket.Conn) *Client {
	client := &Client{
		socket: socket,
		send:   make(chan []byte, MessageBufferSize),
		player: r,
	}
	r.Join(client)
	defer func() { r.leave(client) }()
	go client.Write()
	client.Read()
	return client
}
