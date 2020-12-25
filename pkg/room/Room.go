package room

import (
	"fmt"
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
	for _, user := range meta.Watchers {
		if id == user.ID {
			return true
		}
	}
	return false
}

func (r *Room) Join(usr user.User) {
	meta, _ := r.Store.Find(r.ID)
	watcher, err := meta.FindWatcher(usr.ID)
	if err == nil {
		return
	}
	watcher = NewWatcher(usr)
	watcher.Seek = meta.Seek
	watcher.VideoID = meta.CurrentVideo.ID

	if len(meta.Watchers) == 0 {
		meta.Host = usr.ID
		watcher.IsHost = true
	}

	meta.Watchers = append(meta.Watchers, watcher)
	r.Store.Update(meta)

	r.SendClientEvent(
		Event{
			Action:   EVNT_USER_UPDATE,
			Watchers: meta.Watchers,
		})
}

func (r *Room) SendClientEvent(evt Event) {
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
			r.SendClientEvent(evnt)
		}
	}
}

func (r *Room) SetControls(controls bool) {
	meta, _ := r.Store.Find(r.ID)
	meta.Controls = controls
	r.Store.Update(meta)
}

func (r *Room) AddVideo(video Video) {
	meta, _ := r.Store.Find(r.ID)
	meta.Queue = append(meta.Queue, video)
	r.SetQueue(meta.Queue)
}
func (r *Room) GetVideo() Video {
	meta, _ := r.Store.Find(r.ID)
	return meta.CurrentVideo
}

func (r *Room) GetType() string {
	meta, _ := r.Store.Find(r.ID)
	return meta.Type
}

func (r *Room) SetQueue(queue []Video) bool {
	meta, _ := r.Store.Find(r.ID)

	for i := range queue {
		v := &queue[i]
		if v.ID == "" {
			v.ID = ksuid.New().String()
		}
	}

	meta.Queue = queue
	r.Store.Update(meta)
	r.SendClientEvent(Event{
		Action: EVNT_UPDATE_QUEUE,
		Queue:  meta.Queue,
	})

	if meta.CurrentVideo.ID == "" {
		r.ChangeVideo()
		return false
	}

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

	r.SendClientEvent(Event{
		Action: EVNT_UPDATE_HOST,
		Host:   meta.Host,
	})
}

func (r *Room) GetUser(id string) (RoomWatcher, error) {
	meta, _ := r.Store.Find(r.ID)
	for _, user := range meta.Watchers {
		if user.ID == id {
			return user, nil
		}
	}
	return RoomWatcher{}, fmt.Errorf("User Not found with id: %s", id)
}

func (r *Room) SetPlaying(state bool) {
	meta, _ := r.Store.Find(r.ID)
	meta.Playing = state
	r.Store.Update(meta)
}

func (r *Room) Leave(id string) {
	meta, _ := r.Store.Find(r.ID)
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
	meta, _ := r.Store.Find(r.ID)
	if meta.GetLastVideo().ID == user.VideoID {
		return
	}
	for i := range meta.Watchers {
		u := &meta.Watchers[i]
		if u.Seek < float32(1) && u.ID != user.ID {
			return
		} else if u.ID == user.ID {
			u.Seek = float32(1)
		}
	}
	r.Store.Update(meta)
	r.ChangeVideo()
}

func (r *Room) ChangeVideo() {
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
	})
	r.SendClientEvent(Event{
		Action: EVNT_UPDATE_QUEUE,
		Queue:  meta.Queue,
	})

}

func (r *Room) SeenUser(rw RoomWatcher) {
	meta, _ := r.Store.Find(r.ID)
	meta.UpdateWatcher(rw)

	if meta.Host == rw.ID {
		meta.Seek = rw.Seek
	}

	r.Store.Update(meta)
	r.SendClientEvent(Event{
		Action:   EVT_ON_PROGRESS_UPDATE,
		Watchers: meta.Watchers,
	})
}

// func (r *Room) UpdateUser(u user.User) {
// 	found := false
// 	for i := range meta.Users {
// 		user := &meta.Users[i]

// 		if u.Name == user.Name {
// 			if user.IsHost {
// 				meta.Seek = u.Seek
// 			}
// 			user.LastSeen = time.Now()
// 			user.Seek = u.Seek
// 			user.VideoId = u.VideoId
// 			found = true
// 			room.SendUsersProgress()
// 			return
// 		}

// 	}
// 	if !found {
// 		room.(u)
// 	}
// }
