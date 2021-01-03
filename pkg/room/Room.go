package room

import (
	"fmt"
	"sync"
	"time"
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
	r.mutex.Lock()
	meta, _ := r.Store.Find(r.ID)
	_, err := meta.FindWatcher(id)
	r.mutex.Unlock()
	return err == nil

}

func (r *Room) Join(usr user.User) {
	r.mutex.Lock()
	meta, _ := r.Store.Find(r.ID)
	watcher, err := meta.FindWatcher(usr.ID)
	if err == nil {
		r.mutex.Unlock()
		return
	}
	watcher = NewWatcher(usr)
	watcher.Seek = meta.Seek
	watcher.VideoID = meta.CurrentVideo.ID

	if len(meta.Watchers) == 0 {
		meta.Host = usr.ID
		watcher.IsHost = true
	}

	meta.AddWatcher(watcher)
	r.Store.Update(meta)

	r.mutex.Unlock()
	r.SendClientEvent(
		Event{
			Action:   EVNT_USER_UPDATE,
			Watchers: meta.Watchers,
		})
}

func (r *Room) SendClientEvent(evt Event) {
	if evt.Watcher.ID == "" {
		evt.Watcher = SERVER_USER
	}
	log.Infof("Sending event %s to all clients", evt.Action)
	for client := range r.clients {
		client.send <- evt.ToBytes()
	}
}

func (r *Room) Stop() {
	log.Info("Room Stopping")
	r.Status = ROOM_STATUS_STOPPING
	for client := range r.clients {
		delete(r.clients, client)
	}
	r.Status = ROOM_STATUS_STOPPED
}

func (r *Room) PurgeUsers() bool {
	r.mutex.Lock()
	meta, _ := r.Store.Find(r.ID)
	size := len(meta.Watchers)

	for i := range meta.Watchers {
		wtchr := &meta.Watchers[i]
		if wtchr.LastSeen.Add(10 * time.Second).Before(time.Now()) {
			r.Leave(wtchr.ID)
			size = size - 1
		}
	}
	r.mutex.Unlock()
	return size == 0
}
func (r *Room) DeleteIfEmpty() {
	r.mutex.Lock()
	meta, _ := r.Store.Find(r.ID)
	if meta.Owner == "" {
		log.Infof("No Owner was created annon deleting")
		r.Store.Delete(r.ID)
	}
	r.mutex.Unlock()
}

func (r *Room) Run() {
	r.Status = ROOM_STATUS_RUNNING
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
			evnt, err := processEvent(msg)
			if err != nil {
				return
			}
			log.Infof("sendind message: %v", evnt)
			r.HandleEvent(evnt)
			//r.SendClientEvent(evnt)
		}
	}
}

func (r *Room) SetSettings(settings RoomSettings) {
	r.mutex.Lock()
	meta, _ := r.Store.Find(r.ID)
	meta.Settings = settings
	r.Store.Update(meta)
	r.mutex.Unlock()
}

func (r *Room) AddVideo(video Video, rw RoomWatcher) {
	r.mutex.Lock()
	meta, _ := r.Store.Find(r.ID)
	meta.Queue = append(meta.Queue, video)
	r.SetQueue(meta.Queue, rw)
	r.mutex.Unlock()
}
func (r *Room) GetVideo() Video {
	r.mutex.Lock()
	meta, _ := r.Store.Find(r.ID)
	r.mutex.Unlock()
	return meta.CurrentVideo
}

func (r *Room) GetType() string {
	r.mutex.Lock()
	meta, _ := r.Store.Find(r.ID)
	r.mutex.Unlock()
	return meta.Type
}

func (r *Room) SetQueue(queue []Video, rw RoomWatcher) bool {
	r.mutex.Lock()
	meta, _ := r.Store.Find(r.ID)

	for i := range queue {
		v := &queue[i]
		if v.ID == "" {
			v.ID = ksuid.New().String()
		}
	}
	meta.Queue = queue
	r.Store.Update(meta)
	r.mutex.Unlock()

	if meta.CurrentVideo.ID == "" {
		r.ChangeVideo(rw)
		return false
	}
	r.SendClientEvent(Event{
		Action:  EVNT_UPDATE_QUEUE,
		Queue:   meta.Queue,
		Watcher: rw,
	})
	return true
}

func (r *Room) SetHost(id string) {
	r.mutex.Lock()
	meta, _ := r.Store.Find(r.ID)
	meta.Host = id
	for i := range meta.Watchers {
		if meta.Watchers[i].ID == id {
			meta.Watchers[i].IsHost = true
		}
	}
	r.Store.Update(meta)
	r.mutex.Unlock()
	r.SendClientEvent(Event{
		Action: EVNT_UPDATE_HOST,
		Host:   meta.Host,
	})
}

func (r *Room) GetUser(id string) (RoomWatcher, error) {
	r.mutex.Lock()
	meta, _ := r.Store.Find(r.ID)
	for _, user := range meta.Watchers {
		if user.ID == id {
			r.mutex.Unlock()
			return user, nil
		}
	}
	r.mutex.Unlock()
	return RoomWatcher{}, fmt.Errorf("User Not found with id: %s", id)
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

	r.SendClientEvent(Event{
		Action:   EVNT_USER_UPDATE,
		Watchers: meta.Watchers,
	})
}
func (r *Room) SetSeek(seek float32) {
	meta, _ := r.Store.Find(r.ID)
	meta.Seek = seek
	r.Store.Update(meta)
	r.SendClientEvent(Event{
		Action: EVNT_SEEK_TO_USER,
		Seek:   meta.Seek,
	})
}

func (r *Room) HandleFinish(user RoomWatcher) {
	log.Infof("User %, Has finished! Seek = %f", user.Name, user.Seek)
	r.mutex.Lock()
	user.Seek = float32(1)
	meta, _ := r.Store.Find(r.ID)
	meta.UpdateWatcher(user)

	if !meta.Settings.AutoSkip {
		r.mutex.Unlock()
		return
	}

	if meta.GetLastVideo().ID == user.VideoID {
		r.mutex.Unlock()
		return
	}

	for i := range meta.Watchers {
		u := &meta.Watchers[i]
		if u.Seek < float32(1) && u.ID != user.ID {
			r.mutex.Unlock()
			return
		}
	}
	r.Store.Update(meta)
	r.mutex.Unlock()
	r.ChangeVideo(user)
}

func (r *Room) ChangeVideo(rw RoomWatcher) {
	meta, _ := r.Store.Find(r.ID)
	if len(meta.Queue) == 0 {
		meta.UpdateHistory(meta.CurrentVideo)
		meta.CurrentVideo = Video{}
	} else {
		video := meta.Queue[0]
		meta.Queue = meta.Queue[1:]
		meta.UpdateHistory(meta.CurrentVideo)
		meta.CurrentVideo = video
	}
	r.Store.Update(meta)

	r.SendClientEvent(Event{
		Action:       EVT_VIDEO_CHANGE,
		CurrentVideo: meta.CurrentVideo,
		Watcher:      rw,
	})
	r.SendClientEvent(Event{
		Action:  EVNT_UPDATE_QUEUE,
		Queue:   meta.Queue,
		Watcher: rw,
	})

}

func (r *Room) SeenUser(rw RoomWatcher) {
	meta, _ := r.Store.Find(r.ID)

	err := meta.UpdateWatcher(rw)
	if err != nil {
		meta.AddWatcher(rw)
	}

	if meta.Host == rw.ID {
		meta.Seek = rw.Seek
	}

	r.Store.Update(meta)
	r.SendClientEvent(Event{
		Action:   EVT_ON_PROGRESS_UPDATE,
		Watchers: meta.Watchers,
	})
}
