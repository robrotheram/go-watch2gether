package discordbot

import (
	"fmt"
	"strings"
	"watch2gether/pkg/datastore"
	"watch2gether/pkg/discordbot/command"
	user "watch2gether/pkg/user"

	"github.com/bwmarrin/discordgo"
	log "github.com/sirupsen/logrus"
)

var DiscordUser user.User

type DiscordBot struct {
	*datastore.Datastore
	token   string
	status  string
	session *discordgo.Session
	baseurl string
	voice   *discordgo.VoiceConnection
}

func NewDiscordBot(datastore *datastore.Datastore, token string, baseurl string) (*DiscordBot, error) {
	DiscordUser = user.NewUser("DiscordBot", user.USER_TYPE_DISCORD)
	bot := DiscordBot{
		Datastore: datastore,
		token:     token,
		status:    "INIT",
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

func (db *DiscordBot) Close() {
	db.session.Close()
}

var PREFIX = "!w"

func (db *DiscordBot) MessageCreate(s *discordgo.Session, message *discordgo.MessageCreate) {
	guild, err := db.session.Guild(message.GuildID)
	if err != nil {
		s.ChannelMessageSend(message.ChannelID, "Guild not found")
		return
	}
	channel, err := db.session.Channel(message.ChannelID)
	user := message.Author

	content := message.Content
	if len(content) <= len(PREFIX) {
		return
	}
	if content[:len(PREFIX)] != PREFIX {
		return
	}
	content = content[len(PREFIX):]
	if len(content) < 1 {
		return
	}
	args := strings.Fields(content)
	name := strings.ToLower(args[0])
	ctx := command.CommandCtx{
		Datastore: db.Datastore,
		Session:   s,
		Guild:     guild,
		Channel:   channel,
		User:      user,
		Args:      args[1:],
		BaseURL:   db.baseurl,
	}
	cmd, found := command.Commands[name]
	if !found {
		ctx.Reply(fmt.Sprintf("Error: Command %s not found", name))
		return
	}
	err = cmd.Execute(ctx)
	if err != nil {
		ctx.Reply(fmt.Sprintf("Error: %v", err))
	}
}
