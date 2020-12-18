package pkg

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	log "github.com/sirupsen/logrus"
)

type Hub struct {
	rooms map[string]*room
}

// NewHub Creates a new instaiation of the hub
func NewHub() Hub {
	return Hub{rooms: make(map[string]*room)}
}

// Hub Global var storing all rooms.
// TODO: turn into store backed by DB

/*
DeleteRoom
Sends Delete event across all clients and then closes connection before deleting
*/
func (h *Hub) DeleteRoom(roomID string) {
	if h.rooms[roomID].IsDiscord() {
		h.rooms[roomID].Discord.SendMessage("Room is closeing")
	}

	log.Info("DELETING_ROOM")
	evt := Event{
		Action: "ROOM_EXIT",
	}
	b, _ := json.Marshal(evt)
	h.rooms[roomID].forward <- b
	h.rooms[roomID].quit <- true
	delete(h.rooms, roomID)
}

func (h *Hub) GetRoom(roomid string) (*room, bool) {
	room, found := h.rooms[roomid]
	return room, found
}

func (h *Hub) FindRoom(roomName string) (*room, bool) {
	for _, v := range h.rooms {
		if v.Meta.Name == roomName {
			return v, true
		}
	}
	return nil, false
}

func (h *Hub) NewRoom(room string) *room {
	log.Info("Creaeting New room:" + room)
	return h.AddRoom(newRoom(room))
}

func (h *Hub) AddRoom(room *room) *room {
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
			go r.run()
		}
	}
}

// CleanUP Every 5 seconds go through each room and check to see if there was a delete
func (hub *Hub) CleanUP() {
	log.Info("Staring Cleanup Routine")
	for {
		time.Sleep(5 * time.Second)
		log.Info("Checking Room Infomation")
		for _, room := range hub.rooms {
			room.PurgeUsers()
			if len(room.Meta.Users) == 0 {
				if !room.IsDiscord() {
					go hub.DeleteRoom(room.ID)
				}
			}
		}
	}
}

func (h Hub) DiscordConnect(s *discordgo.Session, e *discordgo.Connect) {
	log.Info("Connected to Something")
	fmt.Println(e)
}

// This function will be called (due to AddHandler above) every time a new
// message is created on any channel that the authenticated bot has access to.
func (h Hub) MessageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {

	// Ignore all messages created by the bot itself
	// This isn't required in this specific example but it's a good practice.
	if m.Author.ID == s.State.User.ID {
		return
	}
	// If the message is "ping" reply with "Pong!"

	args := strings.Fields(m.Content)
	if args[0] != "!w2g" || len(args) < 2 {
		return
	}

	if args[1] == "start" {
		room := newRoom(m.ChannelID)
		room.ID = m.ChannelID
		room.AddDiscord(s, m.ChannelID)
		h.AddRoom(room)
		s.ChannelMessageSend(m.ChannelID, "You room has been created: https://watch2gether.exceptionerror.io/room/"+m.ChannelID)
	}

	if args[1] == "stop" {
		room, _ := h.GetRoom(m.ChannelID)
		h.DeleteRoom(room.ID)
	}

	if args[1] == "add" {
		room, _ := h.GetRoom(m.ChannelID)

		h.DeleteRoom(room.ID)
	}

	// if args[1] == "status" {
	// 	r, _ := h.GetRoom("bob")

	// 	s.ChannelMessageSend(m.ChannelID, "Pong!"+r.Status)
	// }

	// // If the message is "pong" reply with "Ping!"
	// if args[1] == "pong" {
	// 	s.ChannelMessageSend(m.ChannelID, "Ping!")
	// }
}
