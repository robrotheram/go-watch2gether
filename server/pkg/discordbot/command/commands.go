package command

import (
	"sort"
	"watch2gether/pkg/datastore"
	"watch2gether/pkg/room"
	meta "watch2gether/pkg/roomMeta"

	"github.com/bwmarrin/discordgo"
)

type CommandCtx struct {
	*datastore.Datastore
	Session *discordgo.Session
	Guild   *discordgo.Guild
	Channel *discordgo.Channel
	User    *discordgo.User
	Args    []string
	BaseURL string
}

func (ctx *CommandCtx) GetHubRoom() (*room.Room, bool) {
	return ctx.Hub.GetRoom(ctx.Guild.ID)
}

func (ctx *CommandCtx) GetMeta() (*meta.Meta, error) {
	return ctx.Rooms.Find(ctx.Guild.ID)
}

func (ctx *CommandCtx) SaveMeta(meta *meta.Meta) error {
	return ctx.Rooms.Update(meta)
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

func SortKeys(m map[string]Command) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}
