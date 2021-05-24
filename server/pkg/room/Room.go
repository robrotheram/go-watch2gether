package room

import (
	"sync"
	"time"
	"watch2gether/pkg/audioBot"
	events "watch2gether/pkg/events"
	meta "watch2gether/pkg/roomMeta"
	"watch2gether/pkg/user"

	log "github.com/sirupsen/logrus"
)

const ROOM_TYPE_DISCORD = "DISCORD"

type Room struct {
	// forward is a channel that holds incoming messages
	// that should be forwarded to the other clients.
	forward chan []byte
	// join is a channel for clients wishing to join the room.
	join chan *Client
	// leave is a channel for clients wishing to leave the room.
	leave chan *Client
	//Channel to quit the room
	quit    chan bool
	clients map[*Client]bool
	Status  string
	//Discode Bot
	Bot *audioBot.AudioBot

	// clients holds all current clients in this room.
	ID    string
	Store *meta.RoomStore
	sync.Mutex
}

func New(meta *meta.Meta, rs *meta.RoomStore) *Room {
	return &Room{
		forward: make(chan []byte),
		join:    make(chan *Client),
		leave:   make(chan *Client),
		clients: make(map[*Client]bool),
		quit:    make(chan bool),
		Status:  "Initilised",
		ID:      meta.ID,
		Store:   rs,
	}
}

func (r *Room) ContainsUserID(id string) bool {
	meta, _ := r.Store.Find(r.ID)
	_, err := meta.FindWatcher(id)
	return err == nil
}

func (r *Room) RegisterBot(bot *audioBot.AudioBot) error {
	bot.RegisterToRoom(r.forward)
	r.Bot = bot
	return nil //r.bot.Start()

}

func (r *Room) Join(usr user.User) {

	meta, _ := r.Store.Find(r.ID)
	watcher, err := meta.FindWatcher(usr.ID)
	if err == nil {

		return
	}
	watcher = user.NewWatcher(usr)
	watcher.Seek = meta.GetHostSeek()
	watcher.VideoID = meta.CurrentVideo.ID

	if len(meta.Watchers) == 0 {
		meta.Host = usr.ID
		watcher.IsHost = true
	}

	meta.AddWatcher(watcher)
	r.Store.Update(meta)
	r.Send(meta)
}

func (r *Room) UpdateClients() {
	meta, _ := r.Store.Find(r.ID)
	r.Send(meta)
}
func (r *Room) Send(meta *meta.Meta) {
	r.SendStateToClient(events.RoomState{Meta: *meta, Action: events.EVNT_UPDATE_STATE})
}
func (r *Room) SendStateToClient(state events.RoomState) {
	for client := range r.clients {
		client.send <- state.ToBytes()
	}
}

func (r *Room) Stop() {
	log.Info("Room Stopping")
	r.Status = events.ROOM_STATUS_STOPPING
	for client := range r.clients {
		delete(r.clients, client)
	}
	r.Status = events.ROOM_STATUS_STOPPED
}

func (r *Room) PurgeUsers(force bool) bool {
	r.Lock()
	defer r.Unlock()
	meta, err := r.Store.Find(r.ID)
	defer r.Store.Update(meta)
	if err != nil {
		return false
	}

	size := len(meta.Watchers)
	for i := range meta.Watchers {
		wtchr := &meta.Watchers[i]
		if wtchr.Type != user.DISCORD_BOT.Type || force {
			if wtchr.LastSeen.Add(10 * time.Second).Before(time.Now()) {
				meta.RemoveWatcher(wtchr.ID)
				size = size - 1
			}
		}
	}
	return size == 0
}
func (r *Room) DeleteIfEmpty() {
	meta, _ := r.Store.Find(r.ID)
	if meta.Owner == "" {
		log.Infof("No Owner was created annon deleting")
		r.Store.Delete(r.ID)
	}
}

func (r *Room) Run() {
	r.Status = events.ROOM_STATUS_RUNNING
	for {
		select {
		case <-r.quit:
			r.Stop()
			return
		case client := <-r.join:
			r.clients[client] = true
		case client := <-r.leave:
			// leaving
			r.Disconnect(client.user)
		case msg := <-r.forward:
			// forward message to all clients
			evnt, err := events.ProcessEvent(msg)
			if err != nil {
				return
			}
			r.HandleEvent(evnt)
		}
	}
}
func (r *Room) Disconnect(id string) {
	meta, _ := r.Store.Find(r.ID)
	defer r.Store.Update(meta)
	for k := range r.clients {
		if k.user == id {
			log.Debug("user leaving the room")
			meta.RemoveWatcher(id)
			delete(r.clients, k)
			if k.active {
				close(k.send)
				k.active = false
			}
		}
	}
}

func (r *Room) HandleEvent(evt events.Event) {
	r.Lock()
	defer r.Unlock()
	meta, _ := r.Store.Find(r.ID)
	roomState, err := evt.Handle(meta)
	if err != nil {
		log.Warnf("error handling event %v", err)
		return
	}
	r.Store.Update(&roomState.Meta)
	if r.Bot != nil {
		r.Bot.Send(roomState)
	}
	r.SendStateToClient(roomState)
}

func (r *Room) GetType() string {
	meta, _ := r.Store.Find(r.ID)
	return meta.Type
}
