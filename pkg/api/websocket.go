package api

import (
	"encoding/json"
	"net/http"
	"time"
	"w2g/pkg/controllers"
	"w2g/pkg/media"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
)

const WEBPLAYER = controllers.PlayerType("WEB_PLAYER")

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
	user      User
	contoller *controllers.Controller
	socket    *websocket.Conn
	send      chan []byte
	done      chan any
	progress  media.MediaDuration
	running   bool
	exitCode  controllers.PlayerExitCode
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
			log.Debugf("ERROR decoding %v", err)
			c.contoller.RemoveListner(c.id)
			c.contoller.Leave(c.id, c.user.Username)
			return
		}
		var event controllers.Event
		err = json.Unmarshal(msg, &event)
		if err == nil {
			c.UpdateDuration(event.State.Current.Progress)
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

func (wb *Client) Id() string {
	return wb.id
}

func (wb *Client) Meta() controllers.PlayerMeta {
	return controllers.PlayerMeta{
		Id:       wb.id,
		Type:     WEBPLAYER,
		Running:  wb.running,
		Progress: wb.progress,
		User:     wb.user.Username,
	}
}

func (wb *Client) Play(url string, start int) (controllers.PlayerExitCode, error) {
	wb.progress = media.MediaDuration{
		Progress: 0,
	}
	wb.running = true
	<-wb.done
	return wb.exitCode, nil
}

func (wb *Client) Pause() {}

func (wb *Client) Unpause() {
	wb.running = true
}

func (wb *Client) Stop() {
	wb.exitCode = controllers.STOP_EXITCODE
	if wb.running {
		wb.running = false
		wb.done <- "STOP"
	}
}
func (wb *Client) Close() {
	wb.exitCode = controllers.EXIT_EXITCODE
}

func (wb *Client) Status() bool {
	return wb.running
}

func (wb *Client) UpdateDuration(duration media.MediaDuration) {
	wb.progress = duration
}

func (wb *Client) Seek(seconds time.Duration) {
	wb.progress.Progress = seconds
}

func NewClient(socket *websocket.Conn, contoller *controllers.Controller, user User) *Client {
	client := &Client{
		id:        uuid.NewString(),
		user:      user,
		socket:    socket,
		send:      make(chan []byte, MessageBufferSize),
		done:      make(chan any),
		contoller: contoller,
		exitCode:  controllers.STOP_EXITCODE,
	}
	go client.Read()
	go client.Write()
	return client
}
