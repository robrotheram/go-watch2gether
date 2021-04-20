package hub

import (
	"time"
	"watch2gether/pkg/events"
	"watch2gether/pkg/room"
	"watch2gether/pkg/user"

	log "github.com/sirupsen/logrus"
)

type Hub struct {
	Server_ID string
	Rooms     map[string]*room.Room
}

// NewHub Creates a new instaiation of the hub
func NewHub() *Hub {
	return &Hub{Rooms: make(map[string]*room.Room)}
}

// Hub Global var storing all rooms.
// TODO: turn into store backed by DB

/*
DeleteRoom
Sends Delete event across all clients and then closes connection before deleting
*/
func (h *Hub) DeleteRoom(roomID string) {
	// if h.rooms[roomID].IsDiscord() {
	// 	h.rooms[roomID].Discord.SendMessage("Room is closeing")
	// }

	log.Info("DELETING_ROOM")
	h.Rooms[roomID].HandleEvent(events.Event{Action: events.EVT_ROOM_EXIT})

	//h.rooms[roomID].quit <- true
	delete(h.Rooms, roomID)
}

func (h *Hub) GetRoom(roomid string) (*room.Room, bool) {
	room, found := h.Rooms[roomid]
	return room, found
}

func (h *Hub) FindRoom(id string) (*room.Room, bool) {
	for _, v := range h.Rooms {
		if v.ID == id {
			return v, true
		}
	}
	return nil, false
}

// func (h *Hub) NewRoom(meta *room.Meta) *room.Room {
// 	log.Info("Creaeting New room:" + meta.Name)
// 	return h.AddRoom(room.New(meta, h.))
// }

func (h *Hub) AddRoom(room *room.Room) *room.Room {
	log.Info("Adding New room:" + room.ID)
	if _, ok := h.GetRoom(room.ID); ok {
		return nil
	}

	h.Rooms[room.ID] = room
	h.StartRoom(room.ID)
	return room
}

func (h *Hub) StartRoom(roomID string) {
	if r, ok := h.GetRoom(roomID); ok {
		if r.Status != "running" {
			log.Info("Starting Room: " + roomID)
			go r.Run()
		}
	}
}

//CleanUP Every 5 seconds go through each room and check to see if there was a delete

func (hub *Hub) CleanUP(usrStore *user.UserStore) {
	log.Info("Staring Cleanup Routine")
	for {
		time.Sleep(5 * time.Second)
		//log.Info("Checking Room Infomation")
		for _, hubRoom := range hub.Rooms {
			isEmpty := hubRoom.PurgeUsers()
			if isEmpty && hubRoom.GetType() != room.ROOM_TYPE_DISCORD {
				go hub.DeleteRoom(hubRoom.ID)
			}
		}
		annonUsers, _ := usrStore.FindAllByField("Type", user.USER_TYPE_ANON)
		unknownUsers, _ := usrStore.FindAllByField("Type", "")
		users := append(annonUsers, unknownUsers...)
		for _, u := range users {
			found := false
			for _, room := range hub.Rooms {
				if room.ContainsUserID(u.ID) {
					found = true
				}
			}
			if !found {
				log.Infof("Removing user: %s, due to being annon", u.Username)
				go usrStore.Delete(u.ID)
			}
		}
	}
}
