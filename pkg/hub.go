package pkg

import (
	"time"
	"watch2gether/pkg/room"

	log "github.com/sirupsen/logrus"
)

type Hub struct {
	Server_ID string
	rooms     map[string]*room.Room
}

// NewHub Creates a new instaiation of the hub
func NewHub() Hub {
	return Hub{rooms: make(map[string]*room.Room)}
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
	h.rooms[roomID].HandleEvent(room.Event{Action: room.EVT_ROOM_EXIT})

	//h.rooms[roomID].quit <- true
	delete(h.rooms, roomID)
}

func (h *Hub) GetRoom(roomid string) (*room.Room, bool) {
	room, found := h.rooms[roomid]
	return room, found
}

func (h *Hub) FindRoom(id string) (*room.Room, bool) {
	for _, v := range h.rooms {
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

	h.rooms[room.ID] = room
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

func (hub *Hub) CleanUP() {
	log.Info("Staring Cleanup Routine")
	for {
		time.Sleep(5 * time.Second)
		log.Info("Checking Room Infomation")
		for _, room := range hub.rooms {
			isEmpty := room.PurgeUsers()
			if isEmpty && room.GetType() != "DISCORD" {
				go hub.DeleteRoom(room.ID)
			}
		}
	}
}
