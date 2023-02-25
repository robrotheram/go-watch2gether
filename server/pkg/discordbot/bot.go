package discordbot

import (
	"fmt"
	"strconv"
	"watch2gether/pkg/datastore"
	"watch2gether/pkg/discordbot/commands"
	user "watch2gether/pkg/user"

	"github.com/bwmarrin/discordgo"
	log "github.com/sirupsen/logrus"
)

var DiscordUser user.User

type DiscordBot struct {
	*datastore.Datastore
	token    string
	status   string
	session  *discordgo.Session
	baseurl  string
	voice    *discordgo.VoiceConnection
	clientID string
}

func NewDiscordBot(datastore *datastore.Datastore, token string, baseurl string, clientID string) (*DiscordBot, error) {
	DiscordUser = user.NewUser("DiscordBot", user.USER_TYPE_DISCORD)
	bot := DiscordBot{
		Datastore: datastore,
		token:     token,
		status:    "INIT",
		baseurl:   baseurl,
		clientID:  clientID,
	}
	dg, err := discordgo.New("Bot " + token)

	if err != nil {
		return nil, fmt.Errorf("Error Creating Discord Session: %v", err)
	}

	dg.AddHandler(bot.CommandHandler)

	dg.Identify.Intents = discordgo.IntentsGuildMessages | discordgo.IntentsGuilds | discordgo.IntentsGuildMessages | discordgo.IntentsGuildVoiceStates

	// Open a websocket connection to Discord and begin listening.
	bot.session = dg
	return &bot, nil
}

func (db *DiscordBot) RegisterCommands() error {
	cmds, err := db.session.ApplicationCommands(db.session.State.User.ID, "")
	doesCmdExist := func(name string) bool {
		for _, v := range cmds {
			if v.Name == name {
				return true
			}
		}
		return false
	}
	if err != nil {
		log.Warnf("error getting commands: %v", err)
		return err
	}
	for _, v := range cmds {
		cmd, err := commands.Commands.GetCommand(v.Name)
		if err == nil {
			_, err = db.session.ApplicationCommandEdit(db.session.State.User.ID, "", v.ID, &cmd.ApplicationCommand)
			log.Infof("updating command: %s", cmd.Name)
		} else {
			err = db.session.ApplicationCommandDelete(db.session.State.User.ID, "", v.ID)
			log.Infof("removing command: %s", v.Name)
		}
		if err != nil {
			log.Warnf("error updating command: %v", err)
		}
	}
	for _, cmd := range commands.Commands.Cmds {
		if !doesCmdExist(cmd.Name) {
			acc, err := db.session.ApplicationCommandCreate(db.session.State.User.ID, "", &cmd.ApplicationCommand)
			log.Infof("creating command: %s", acc.Name)
			if err != nil {
				log.Warnf("error updating command: %v", err)
			}
		}
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

	//go db.RegisterCommands()
	return nil
}

func (db *DiscordBot) Close() {
	db.session.Close()
}

func (db *DiscordBot) CommandHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {

	guild, err := db.session.Guild(i.GuildID)
	log.Infof("Guild %s -running command: %s", guild.Name, i.ApplicationCommandData().Name)

	if err != nil {
		s.ChannelMessageSend(i.ChannelID, "Guild not found")
		return
	}
	channel, err := db.session.Channel(i.ChannelID)
	user := i.Interaction.Member

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
		Datastore: db.Datastore,
		Session:   s,
		Guild:     guild,
		Channel:   channel,
		User:      user,
		Args:      args,
		BaseURL:   db.baseurl,
	}

	cmd, err := commands.Commands.GetCommand(i.ApplicationCommandData().Name)
	if err != nil {
		ctx.Reply(fmt.Sprintf("%v", err))
		return
	}
	if resp := cmd.Function(ctx); resp != nil {
		s.InteractionRespond(i.Interaction, resp)
	}
}
