package discordbot

import (
	"fmt"
	"strconv"
	"watch2gether/pkg/players"
	"watch2gether/pkg/playlists"

	"github.com/bwmarrin/discordgo"
	"github.com/labstack/gommon/log"
)

func RegisterCommands(session *discordgo.Session) error {
	cmds, err := session.ApplicationCommands(session.State.User.ID, "")
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
		cmd, err := Commands.GetCommand(v.Name)
		if err == nil {
			_, err = session.ApplicationCommandEdit(session.State.User.ID, "", v.ID, &cmd.ApplicationCommand)
			log.Infof("updating command: %s", cmd.Name)
		} else {
			err = session.ApplicationCommandDelete(session.State.User.ID, "", v.ID)
			log.Infof("removing command: %s", v.Name)
		}
		if err != nil {
			log.Warnf("error updating command: %v", err)
		}
	}
	for _, cmd := range Commands.Cmds {
		if !doesCmdExist(cmd.Name) {
			acc, err := session.ApplicationCommandCreate(session.State.User.ID, "", &cmd.ApplicationCommand)
			log.Infof("creating command: %s", acc.Name)
			if err != nil {
				log.Warnf("error updating command: %v", err)
			}
		}
	}
	log.Info("Updating Commands complete")
	return nil
}

func RegisterCommandHandler(store *players.Store, playlistStore *playlists.PlaylistStore) func(session *discordgo.Session, i *discordgo.InteractionCreate) {

	return func(session *discordgo.Session, i *discordgo.InteractionCreate) {
		guild, err := session.Guild(i.GuildID)
		log.Infof("Guild %s -running command: %s", guild.Name, i.ApplicationCommandData().Name)

		if err != nil {
			session.ChannelMessageSend(i.ChannelID, "Guild not found")
			return
		}
		channel, err := session.Channel(i.ChannelID)
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
				for _, a := range arg.Options {
					switch a.Type {
					case discordgo.ApplicationCommandOptionString:
						args = append(args, a.StringValue())
					case discordgo.ApplicationCommandOptionInteger:
						args = append(args, strconv.FormatInt(a.IntValue(), 10))
					}
				}
			}
		}

		ctx := CommandCtx{
			Session:   session,
			Guild:     guild,
			Channel:   channel,
			User:      user,
			Args:      args,
			Store:     store,
			Playlists: playlistStore,
		}
		cmd, err := Commands.GetCommand(i.ApplicationCommandData().Name)
		if err != nil {
			ctx.Reply(fmt.Sprintf("%v", err))
			return
		}
		if resp := cmd.Function(ctx); resp != nil {
			session.InteractionRespond(i.Interaction, resp)
		}
	}
}
