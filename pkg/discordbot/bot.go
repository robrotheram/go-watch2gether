package discord

import (
	"fmt"
	"net/url"
	"strings"
	"watch2gether/pkg/hub"
	"watch2gether/pkg/media"
	"watch2gether/pkg/room"
	"watch2gether/pkg/roombot"
	user "watch2gether/pkg/user"
	"watch2gether/pkg/utils"

	"github.com/bwmarrin/discordgo"
	"github.com/segmentio/ksuid"
	log "github.com/sirupsen/logrus"
)

var DiscordUser user.User

type DiscordBot struct {
	hub       *hub.Hub
	token     string
	status    string
	session   *discordgo.Session
	roomStore *room.RoomStore
	playlists *media.PlayistStore
	baseurl   string
	voice     *discordgo.VoiceConnection
}

func NewDiscordBot(h *hub.Hub, roomStore *room.RoomStore, playlists *media.PlayistStore, token string, baseurl string) (*DiscordBot, error) {
	DiscordUser = user.NewUser("DiscordBot", user.USER_TYPE_DISCORD)
	bot := DiscordBot{
		hub:       h,
		token:     token,
		status:    "INIT",
		roomStore: roomStore,
		playlists: playlists,
		baseurl:   baseurl,
	}
	dg, err := discordgo.New("Bot " + token)
	if err != nil {
		return nil, fmt.Errorf("Error Creating Discord Session: %v", err)
	}

	dg.AddHandler(bot.MessageCreate)

	dg.Identify.Intents = discordgo.IntentsGuildMessages | discordgo.IntentsGuilds | discordgo.IntentsGuildMessages | discordgo.IntentsGuildVoiceStates

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

	fmt.Println(db.session.State)
	return nil
}

func (db *DiscordBot) GetUserVoiceChannel(user string) (string, error) {
	for _, g := range db.session.State.Guilds {
		for _, v := range g.VoiceStates {
			if v.UserID == user {
				return v.ChannelID, nil
			}
		}
	}
	return "", fmt.Errorf("Channel Not found")
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
		meta := room.NewMeta(guild.Name, user.User{ID: m.Author.ID, Type: user.USER_TYPE_BASIC})
		meta.ID = m.ChannelID
		meta.Type = "DISCORD"
		db.roomStore.Create(meta)
		r := room.New(meta, db.roomStore)
		//room.AddDiscord(s, m.ChannelID)
		db.hub.AddRoom(r)
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("You room has been created: https://%s/room/%s", db.baseurl, m.ChannelID))
		return
	}
	if args[1] == "load" {
		r, ok := db.hub.GetRoom(m.GuildID)
		if !ok {
			s.ChannelMessageSend(m.ChannelID, "Room not found")
			return
		}
		playlists, err := db.playlists.FindByField("RoomID", m.GuildID)
		if err != nil {
			s.ChannelMessageSend(m.ChannelID, "Playlists not found")
			return
		}

		playlistName := strings.TrimSuffix(strings.Join(args[2:], " "), " ")
		fmt.Printf("Searching for playlist %s: \n", playlistName)
		for _, playlist := range playlists {
			if playlist.Name == playlistName {
				queue := r.GetQueue()
				queue = append(queue, playlist.Videos...)
				r.SetQueue(queue, user.DISCORD_BOT)
				return
			}
		}
	}

	if args[1] == "join" {

		r, ok := db.hub.GetRoom(m.GuildID)
		if !ok {
			s.ChannelMessageSend(m.ChannelID, "Room not found")
			return
		}
		vc, err := db.GetUserVoiceChannel(m.Author.ID)
		if err != nil {
			s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("User not connected to voice channel"))
		}
		voice, err := db.session.ChannelVoiceJoin(m.GuildID, vc, false, true)
		if err != nil {
			s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("User not connected to voice channel"))
		}
		bot := roombot.NewAudioBot("", m.ChannelID, voice, s)
		err = bot.Start()
		if err != nil {
			s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Bot error %w", err))
		}
		r.RegisterBot(bot)

		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Guild ID: %s, Channel ID: %s", m.GuildID, vc))
		db.voice, err = db.session.ChannelVoiceJoin(m.GuildID, vc, false, true)

		return
	}

	if args[1] == "stop" {
		room, _ := db.hub.GetRoom(m.GuildID)
		db.hub.DeleteRoom(room.ID)
		return
	}

	if args[1] == "add" && len(args) == 3 {
		u, err := url.ParseRequestURI(args[2])
		if err != nil {
			s.ChannelMessageSend(m.ChannelID, "Not a valid URL")
			return
		}
		r, ok := db.hub.GetRoom(m.GuildID)
		if !ok {
			s.ChannelMessageSend(m.ChannelID, "Room not found: "+m.GuildID)
			return
		}
		document, err := utils.Scrape(u.String(), 1)
		if err != nil {
			log.Error(err)
			s.ChannelMessageSend(m.ChannelID, "Video Error not found")
			return
		}

		video := media.Video{ID: ksuid.New().String(), Title: document.Preview.Title, Url: u.String(), User: DiscordUser.Username}
		r.AddVideo(video, user.DISCORD_BOT)

		log.Infof("Vidoe Envent sent : %v", video)
		s.ChannelMessageSend(m.ChannelID, "Video Added ID:"+video.ID)
		return
	}

	if args[1] == "skip" {
		r, ok := db.hub.GetRoom(m.GuildID)
		if !ok {
			s.ChannelMessageSend(m.ChannelID, "Room not found")
			return
		}
		r.ChangeVideo(user.DISCORD_BOT)
		s.ChannelMessageSend(m.ChannelID, "Video Skiped")
		return
	}

	if args[1] == "status" {
		r, ok := db.hub.GetRoom(m.GuildID)
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
