package room

import (
	"net/http"

	"github.com/gorilla/websocket"
)

// client represents a single chatting user.
type Client struct {
	// socket is the web socket for this client.
	socket *websocket.Conn
	// send is a channel on which messages are sent.
	send chan []byte
	// room is the room this client is chatting in.
	room   *Room
	active bool
	user   string
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
			return
		}
		c.room.forward <- msg
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

func NewClient(r *Room, socket *websocket.Conn, token string) *Client {
	client := &Client{
		socket: socket,
		send:   make(chan []byte, MessageBufferSize),
		room:   r,
		user:   token,
		active: true,
	}
	r.join <- client
	defer func() { r.leave <- client }()
	go client.Write()
	client.Read()
	return client
}
