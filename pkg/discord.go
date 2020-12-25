package pkg

import (
	"fmt"
	"net/url"
	"strings"
	"watch2gether/pkg/room"
	user "watch2gether/pkg/user"

	"github.com/bwmarrin/discordgo"
	"github.com/segmentio/ksuid"
	log "github.com/sirupsen/logrus"
)

var DiscordUser = user.NewUser("DiscordBot")

type DiscordBot struct {
	hub       *Hub
	token     string
	status    string
	session   *discordgo.Session
	roomStore *room.RoomStore
}

func NewDiscordBot(h *Hub, roomStore *room.RoomStore, token string) (*DiscordBot, error) {
	bot := DiscordBot{
		hub:       h,
		token:     token,
		status:    "INIT",
		roomStore: roomStore,
	}
	dg, err := discordgo.New("Bot " + token)
	if err != nil {
		return nil, fmt.Errorf("Error Creating Discord Session: %v", err)
	}

	dg.AddHandler(bot.MessageCreate)

	dg.Identify.Intents = discordgo.MakeIntent(discordgo.IntentsGuildMessages)

	// Open a websocket connection to Discord and begin listening.
	bot.session = dg
	return &bot, nil
}

func (db *DiscordBot) Start() error {
	// Open a websocket connection to Discord and begin listening.
	log.Info("Discord Bot Starting")
	err := db.session.Open()
	if err != nil {
		return fmt.Errorf("Error opening connection: %v", err)
	}
	return nil
}

func (db *DiscordBot) Close() {
	db.session.Close()
}

func (db *DiscordBot) MessageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {

	if m.Author.ID == s.State.User.ID {
		return
	}
	guild, err := db.session.Guild(m.GuildID)
	if err != nil {
		return
	}

	args := strings.Fields(m.Content)
	if args[0] != "!w2g" || len(args) < 2 {
		log.Warn("Discord command not reconsied")
		return
	}

	if args[1] == "start" {
		meta := room.NewMeta(guild.Name, s.State.User.ID)
		meta.ID = m.ChannelID
		meta.Type = "DISCORD"
		db.roomStore.Create(meta)
		r := room.New(meta, db.roomStore)
		//room.AddDiscord(s, m.ChannelID)
		db.hub.AddRoom(r)
		s.ChannelMessageSend(m.ChannelID, "You room has been created: https://watch2gether.exceptionerror.io/room/"+m.ChannelID)
		return
	}

	if args[1] == "stop" {
		room, _ := db.hub.GetRoom(m.ChannelID)
		db.hub.DeleteRoom(room.ID)
		return
	}

	if args[1] == "add" && len(args) == 3 {
		u, err := url.ParseRequestURI(args[2])
		if err != nil {
			s.ChannelMessageSend(m.ChannelID, "Not a valid URL")
			return
		}
		r, ok := db.hub.GetRoom(m.ChannelID)
		if !ok {
			s.ChannelMessageSend(m.ChannelID, "Room not found")
			return
		}
		document, err := Scrape(u.String(), 1)
		if err != nil {
			log.Error(err)
			s.ChannelMessageSend(m.ChannelID, "Video Error not found")
			return
		}

		video := room.Video{ID: ksuid.New().String(), Title: document.Preview.Title, Url: u.String(), User: DiscordUser.Name}
		r.AddVideo(video)

		log.Infof("Vidoe Envent sent : %v", video)
		s.ChannelMessageSend(m.ChannelID, "Video Added ID:"+video.ID)
		return
	}

	if args[1] == "skip" {
		r, ok := db.hub.GetRoom(m.ChannelID)
		if !ok {
			s.ChannelMessageSend(m.ChannelID, "Room not found")
			return
		}
		r.ChangeVideo()
		s.ChannelMessageSend(m.ChannelID, "Video Skiped")
		return
	}

	if args[1] == "status" {
		r, ok := db.hub.GetRoom(m.ChannelID)
		if !ok {
			s.ChannelMessageSend(m.ChannelID, "Room not found")
			return
		}
		vidoe := r.GetVideo()
		if vidoe.ID == "" {
			s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("No Video Playing"))
			return
		}
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Currently Playing: %s", vidoe.Title))
		return
	}

	s.ChannelMessageSend(m.ChannelID, "Command Not reconsided:"+m.Content)
}
