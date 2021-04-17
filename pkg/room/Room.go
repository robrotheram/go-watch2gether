package room

import (
	"fmt"
	"sync"
	"time"
	events "watch2gether/pkg/events"
	"watch2gether/pkg/media"
	"watch2gether/pkg/roombot"
	"watch2gether/pkg/user"

	"github.com/segmentio/ksuid"
	log "github.com/sirupsen/logrus"
)

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
	bot *roombot.AudioBot

	// clients holds all current clients in this room.
	ID    string
	Store *RoomStore
	mutex sync.Mutex
}

func New(meta *Meta, rs *RoomStore) *Room {
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

func (r *Room) RegisterBot(bot *roombot.AudioBot) {
	r.bot = bot
}

func (r *Room) Join(usr user.User) {

	meta, _ := r.Store.Find(r.ID)
	watcher, err := meta.FindWatcher(usr.ID)
	if err == nil {

		return
	}
	watcher = user.NewWatcher(usr)
	watcher.Seek = meta.Seek
	watcher.VideoID = meta.CurrentVideo.ID

	if len(meta.Watchers) == 0 {
		meta.Host = usr.ID
		watcher.IsHost = true
	}

	meta.AddWatcher(watcher)
	r.Store.Update(meta)

	r.SendClientEvent(
		events.Event{
			Action:   events.EVNT_USER_UPDATE,
			Watchers: meta.Watchers,
		})
}

func (r *Room) SendClientEvent(evt events.Event) {
	if evt.Watcher.ID == "" {
		evt.Watcher = user.SERVER_USER
	}
	log.Infof("Sending event %s to all clients", evt.Action)
	for client := range r.clients {
		client.send <- evt.ToBytes()
	}
	if r.bot != nil {
		r.bot.Send(evt)
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

func (r *Room) PurgeUsers() bool {

	meta, _ := r.Store.Find(r.ID)
	size := len(meta.Watchers)

	for i := range meta.Watchers {
		wtchr := &meta.Watchers[i]
		if wtchr.LastSeen.Add(10 * time.Second).Before(time.Now()) {
			r.Leave(wtchr.ID)
			size = size - 1
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

func (r *Room) HandleEvent(evt events.Event) {
	if evt.Watcher.ID == user.SERVER_USER.ID {
		return
	}
	switch evt.Action {
	case events.EVNT_PLAYING:
		r.SetPlaying(true)
		r.SendClientEvent(evt)
	case events.EVNT_PAUSING:
		r.SetPlaying(false)
		r.SendClientEvent(evt)
	case events.EVNT_UPDATE_HOST:
		r.SetHost(evt.Host)
	case events.EVNT_NEXT_VIDEO:
		r.ChangeVideo(evt.Watcher)
	case events.EVNT_SEEK:
		r.SetSeek(evt.Seek)
	case events.EVNT_UPDATE_SETTINGS:
		r.SetSettings(evt.Settings)
	case events.EVNT_SEEK_TO_ME:
		r.SetSeek(evt.Watcher.Seek)
	case events.EVNT_UPDATE_QUEUE:
		r.SetQueue(evt.Queue, evt.Watcher)
	case events.ENVT_FINSH:
		r.HandleFinish(evt.Watcher)
	case events.EVNT_USER_UPDATE:
		r.SeenUser(evt.Watcher)
	case events.EVT_ROOM_EXIT:
		r.DeleteIfEmpty()
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
			r.Leave(client.user)
			delete(r.clients, client)
			close(client.send)
		case msg := <-r.forward:
			// forward message to all clients
			evnt, err := events.ProcessEvent(msg)
			if err != nil {
				return
			}
			log.Infof("sendind message: %v", evnt)
			r.HandleEvent(evnt)
			//r.SendClientEvent(evnt)
		}
	}
}

func (r *Room) SetSettings(settings events.RoomSettings) {

	meta, _ := r.Store.Find(r.ID)
	meta.Settings = settings
	r.Store.Update(meta)

}

func (r *Room) AddVideo(video media.Video, rw user.Watcher) {

	meta, _ := r.Store.Find(r.ID)
	meta.Queue = append(meta.Queue, video)
	r.SetQueue(meta.Queue, rw)

}
func (r *Room) GetVideo() media.Video {

	meta, _ := r.Store.Find(r.ID)

	return meta.CurrentVideo
}

func (r *Room) GetType() string {

	meta, _ := r.Store.Find(r.ID)

	return meta.Type
}

func (r *Room) GetQueue() []media.Video {
	meta, _ := r.Store.Find(r.ID)
	return meta.Queue
}

func (r *Room) SetQueue(queue []media.Video, rw user.Watcher) bool {

	meta, _ := r.Store.Find(r.ID)

	for i := range queue {
		v := &queue[i]
		if v.ID == "" {
			v.ID = ksuid.New().String()
		}
	}
	meta.Queue = queue
	r.Store.Update(meta)

	if meta.CurrentVideo.ID == "" {
		r.ChangeVideo(rw)
		return false
	}
	r.SendClientEvent(events.Event{
		Action:  events.EVNT_UPDATE_QUEUE,
		Queue:   meta.Queue,
		Watcher: rw,
	})
	return true
}

func (r *Room) SetHost(id string) {

	meta, _ := r.Store.Find(r.ID)
	meta.Host = id
	for i := range meta.Watchers {
		if meta.Watchers[i].ID == id {
			meta.Watchers[i].IsHost = true
		}
	}
	r.Store.Update(meta)

	r.SendClientEvent(events.Event{
		Action: events.EVNT_UPDATE_HOST,
		Host:   meta.Host,
	})
}

func (r *Room) GetUser(id string) (user.Watcher, error) {
	meta, _ := r.Store.Find(r.ID)
	for _, user := range meta.Watchers {
		if user.ID == id {
			return user, nil
		}
	}
	return user.Watcher{}, fmt.Errorf("User Not found with id: %s", id)
}

func (r *Room) SetPlaying(state bool) {
	meta, _ := r.Store.Find(r.ID)
	meta.Playing = state
	r.Store.Update(meta)
}

func (r *Room) Leave(id string) {
	meta, _ := r.Store.Find(r.ID)
	if meta == nil {
		return
	}
	meta.RemoveWatcher(id)
	r.Store.Update(meta)

	if meta.Host == id && len(meta.Watchers) > 0 {
		r.SetHost(meta.Watchers[0].ID)
	}

	log.Infof("User: %s Has left the room: %s", id, meta.Name)

	r.SendClientEvent(events.Event{
		Action:   events.EVNT_USER_UPDATE,
		Watchers: meta.Watchers,
	})
}
func (r *Room) SetSeek(seek float32) {
	meta, _ := r.Store.Find(r.ID)
	meta.Seek = seek
	r.Store.Update(meta)
	r.SendClientEvent(events.Event{
		Action: events.EVNT_SEEK_TO_USER,
		Seek:   meta.Seek,
	})
}

func (r *Room) HandleFinish(user user.Watcher) {
	log.Infof("User %, Has finished! Seek = %f", user.Username, user.Seek)

	user.Seek = float32(1)
	meta, _ := r.Store.Find(r.ID)
	meta.UpdateWatcher(user)

	if !meta.Settings.AutoSkip {

		return
	}

	if meta.GetLastVideo().ID == user.VideoID {

		return
	}

	for i := range meta.Watchers {
		u := &meta.Watchers[i]
		if u.Seek < float32(1) && u.ID != user.ID {
			return
		}
	}
	r.Store.Update(meta)

	r.ChangeVideo(user)
}

func (r *Room) ChangeVideo(rw user.Watcher) {
	meta, _ := r.Store.Find(r.ID)
	if len(meta.Queue) == 0 {
		meta.UpdateHistory(meta.CurrentVideo)
		meta.CurrentVideo = media.Video{}
	} else {
		video := meta.Queue[0]
		meta.Queue = meta.Queue[1:]
		meta.UpdateHistory(meta.CurrentVideo)
		meta.CurrentVideo = video
	}
	r.Store.Update(meta)

	r.SendClientEvent(events.Event{
		Action:       events.EVT_VIDEO_CHANGE,
		CurrentVideo: meta.CurrentVideo,
		Watcher:      rw,
	})
	r.SendClientEvent(events.Event{
		Action:  events.EVNT_UPDATE_QUEUE,
		Queue:   meta.Queue,
		Watcher: rw,
	})

}

func (r *Room) SeenUser(rw user.Watcher) {

	meta, _ := r.Store.Find(r.ID)

	err := meta.UpdateWatcher(rw)
	if err != nil {
		meta.AddWatcher(rw)
	}

	if meta.Host == rw.ID {
		meta.Seek = rw.Seek
	}

	r.Store.Update(meta)

	r.SendClientEvent(events.Event{
		Action:   events.EVT_ON_PROGRESS_UPDATE,
		Watchers: meta.Watchers,
	})
}
