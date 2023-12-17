package api

import (
	"encoding/json"
	"net/http"
	"w2g/pkg/controllers"

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
	// socket is the web socket for this client.
	socket *websocket.Conn
	// send is a channel on which messages are sent.
	send chan []byte
	// room is the room this client is chatting in.
	active bool
}

const (
	SocketBufferSize  = 1024
	MessageBufferSize = 256
)

// func (c *Client) Read() {
// 	defer c.socket.Close()
// 	for {
// 		// _, msg, err := c.socket.ReadMessage()
// 		// if err != nil {
// 		// 	fmt.Printf("ERROR decoding %v", err)
// 		// 	return
// 		// }
// 		// c.room.forward <- msg
// 	}
// }

func (c *Client) Write() {
	defer c.socket.Close()
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

func NewClient(socket *websocket.Conn, contoller *controllers.Controller) *Client {
	client := &Client{
		socket: socket,
		send:   make(chan []byte, MessageBufferSize),
		active: true,
	}
	defer func() { contoller.RemoveListner(client) }()
	go client.Write()
	return client
}
