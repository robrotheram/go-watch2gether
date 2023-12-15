package discord

import (
	"fmt"
	"strconv"
	"w2g/pkg/controllers"
	"w2g/pkg/discord/commands"
	"w2g/pkg/discord/components"
	"w2g/pkg/discord/session"
	"w2g/pkg/utils"

	"github.com/bwmarrin/discordgo"
	log "github.com/sirupsen/logrus"
)

type DiscordBot struct {
	token    string
	status   string
	session  *discordgo.Session
	channels *controllers.Hub
	clientID string
	users    map[string]*session.UserSession
}

func NewDiscordBot(config utils.Config) (*DiscordBot, error) {
	bot := DiscordBot{
		token:    config.DiscordToken,
		status:   "INIT",
		clientID: config.DiscordClientID,
		channels: controllers.NewHub(),
		users:    make(map[string]*session.UserSession),
	}
	dg, err := discordgo.New("Bot " + bot.token)

	if err != nil {
		return nil, fmt.Errorf("error creating discord session: %v", err)
	}

	dg.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) {
		log.Println("Bot is up!")
	})

	dg.AddHandler(bot.CommandHandler)

	dg.Identify.Intents = discordgo.IntentsGuildMessages | discordgo.IntentsGuilds | discordgo.IntentsGuildMessages | discordgo.IntentsGuildVoiceStates

	// Open a websocket connection to Discord and begin listening.
	bot.session = dg
	return &bot, nil
}

func (db *DiscordBot) RegisterCommands() error {
	cmds, _ := db.session.ApplicationCommands(db.clientID, "")

	for _, v := range cmds {
		err := db.session.ApplicationCommandDelete(db.clientID, "", v.ID)
		if err != nil {
			log.Warnf("error removing command: %v", err)
		}
		log.Infof("removing command: %s", v.Name)
	}

	for name, cmd := range commands.GetCommands() {
		acc, err := db.session.ApplicationCommandCreate(db.clientID, "", &cmd.ApplicationCommand)
		if err != nil {
			log.Warnf("error updating command %s: %v", name, err)
		}
		log.Infof("creating command: %s", acc.Name)
	}
	log.Info("Updating Commands complete")
	return nil
}

func (db *DiscordBot) Start() error {
	// Open a websocket connection to Discord and begin listening.
	log.Info("Discord Bot Starting")

	err := db.session.Open()
	if err != nil {
		return fmt.Errorf("error opening connection: %v", err)
	}

	db.RegisterCommands()
	return nil
}

func (db *DiscordBot) Close() {
	db.session.Close()
}

func (db *DiscordBot) handleApplicationCommand(s *discordgo.Session, i *discordgo.InteractionCreate) {
	guild, _ := db.session.Guild(i.GuildID)
	channel, _ := db.session.Channel(i.ChannelID)
	user := i.Interaction.Member
	controller := db.channels.Get(i.GuildID)

	args := []string{}
	for _, arg := range i.ApplicationCommandData().Options {
		switch arg.Type {
		case discordgo.ApplicationCommandOptionString:
			args = append(args, arg.StringValue())
		case discordgo.ApplicationCommandOptionInteger:
			args = append(args, strconv.FormatInt(arg.IntValue(), 10))
		case discordgo.ApplicationCommandOptionSubCommand:
			args = append(args, arg.Name)
		}
	}
	ctx := commands.CommandCtx{
		Session:     s,
		Guild:       guild,
		Channel:     channel,
		User:        user,
		Args:        args,
		Controller:  controller,
		UserSession: db.getUserSession(user),
	}
	cmd, err := commands.GetCommand(i.ApplicationCommandData().Name)
	if err != nil {
		ctx.Reply(fmt.Sprintf("%v", err))
		return
	}
	if resp := cmd.Function(ctx); resp != nil {
		err := s.InteractionRespond(i.Interaction, resp)
		if err != nil {
			log.Println(err)
		}
	}
}
func (db *DiscordBot) getUserSession(user *discordgo.Member) *session.UserSession {
	if _, ok := db.users[user.User.ID]; !ok {
		db.users[user.User.ID] = &session.UserSession{}
	}
	return db.users[user.User.ID]
}

func (db *DiscordBot) handleMessageCommand(s *discordgo.Session, i *discordgo.InteractionCreate) {
	guild, _ := db.session.Guild(i.GuildID)
	channel, _ := db.session.Channel(i.ChannelID)
	user := i.Interaction.Member
	controller := db.channels.Get(i.GuildID)
	ctx := components.HandlerCtx{
		Session:     s,
		Guild:       guild,
		Channel:     channel,
		User:        user,
		Controller:  controller,
		UserSession: db.getUserSession(user),
	}
	cmd, err := components.GetHandler(i.MessageComponentData().CustomID)
	if err != nil {
		return
	}

	if resp := cmd.Function(ctx); resp != nil {
		err := s.InteractionRespond(i.Interaction, resp)
		if err != nil {
			log.Println(err)
		}
	}
}

func (db *DiscordBot) CommandHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	switch i.Type {
	case discordgo.InteractionApplicationCommand:
		db.handleApplicationCommand(s, i)
	case discordgo.InteractionMessageComponent:
		db.handleMessageCommand(s, i)
	}
}
