package commands

import (
	"fmt"
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
	User    *discordgo.Member
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

// func (ctx *CommandCtx) Reply(message string) error {
// 	_, err := ctx.Session.ChannelMessageSend(
// 		ctx.Channel.ID,
// 		message,
// 	)
// 	return err
// }
func (ctx *CommandCtx) ReplyEmbed(message *EmbededMessage) error {
	_, err := ctx.Session.ChannelMessageSendEmbed(
		ctx.Channel.ID,
		&message.MessageEmbed,
	)
	return err
}

func (ctx *CommandCtx) Reply(message string) *discordgo.InteractionResponse {
	return &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: message,
		},
	}
}

func (ctx *CommandCtx) Errorf(format string, a ...interface{}) *discordgo.InteractionResponse {
	return &discordgo.InteractionResponse{
		Type: discordgo.InteractionApplicationCommandAutocompleteResult,
		Data: &discordgo.InteractionResponseData{
			Content: fmt.Sprintf(format, a...),
		},
	}
}
