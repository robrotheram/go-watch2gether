package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"w2g/pkg/controllers"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

var Upgrader = &websocket.Upgrader{
	ReadBufferSize:  SocketBufferSize,
	WriteBufferSize: SocketBufferSize,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// client represents a single chatting user.
type Client struct {
	id        string
	contoller *controllers.Controller
	player    *WebPlayer
	// socket is the web socket for this client.
	socket *websocket.Conn
	// send is a channel on which messages are sent.
	send chan []byte
}

const (
	SocketBufferSize  = 1024
	MessageBufferSize = 256
)

func (c *Client) Read() {
	defer c.socket.Close()
	for {
		_, msg, err := c.socket.ReadMessage()
		if err != nil {
			fmt.Printf("ERROR decoding %v", err)
			return
		}
		var event controllers.Event
		err = json.Unmarshal(msg, &event)
		if err == nil {
			c.player.UpdateDuration(event.State.Current.Progress)
		}
	}
}

func (c *Client) Write() {
	defer c.socket.Close()
	defer c.contoller.RemoveListner(c.id)
	for msg := range c.send {
		err := c.socket.WriteMessage(websocket.TextMessage, msg)
		if err != nil {
			return
		}
	}
}

func (c *Client) Send(event controllers.Event) {
	data, err := json.Marshal(event)
	if err == nil {
		c.send <- data
	}
}

func NewClient(socket *websocket.Conn, contoller *controllers.Controller, player *WebPlayer) *Client {
	client := &Client{
		id:        uuid.NewString(),
		socket:    socket,
		send:      make(chan []byte, MessageBufferSize),
		contoller: contoller,
		player:    player,
	}
	go client.Read()
	go client.Write()
	return client
}
