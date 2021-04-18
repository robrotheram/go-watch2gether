package command

import (
	"watch2gether/pkg/datastore"
	"watch2gether/pkg/room"

	"github.com/bwmarrin/discordgo"
)

type CommandCtx struct {
	*datastore.Datastore
	Session *discordgo.Session
	Guild   *discordgo.Guild
	Channel *discordgo.Channel
	User    *discordgo.User
	Args    []string
}

func (ctx *CommandCtx) GetHubRoom() (*room.Room, bool) {
	return ctx.Hub.GetRoom(ctx.Guild.ID)
}

func (ctx *CommandCtx) Reply(message string) error {
	_, err := ctx.Session.ChannelMessageSend(
		ctx.Channel.ID,
		message,
	)
	return err
}

type Command interface {
	GetHelp() string
	Execute(CommandCtx) error
}

type BaseCommand struct {
	Help string
}

func (cmd *BaseCommand) GetHelp() string {
	return cmd.Help
}

var Commands = make(map[string]Command)
