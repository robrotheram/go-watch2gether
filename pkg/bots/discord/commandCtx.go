package discordbot

import (
	"fmt"
	"watch2gether/pkg/channels"
	"watch2gether/pkg/playlists"

	"github.com/bwmarrin/discordgo"
)

type CommandCtx struct {
	*channels.Store
	Playlists *playlists.PlaylistStore
	Session   *discordgo.Session
	Guild     *discordgo.Guild
	Channel   *discordgo.Channel
	User      *discordgo.Member
	Args      []string
	BaseURL   string
}

func (ctx *CommandCtx) ReplyEmbed(message *EmbededMessage) error {
	_, err := ctx.Session.ChannelMessageSendEmbed(
		ctx.Channel.ID,
		&message.MessageEmbed,
	)
	return err
}

func (ctx *CommandCtx) CmdReplyEmbed(message *EmbededMessage) *discordgo.InteractionResponse {
	embeds := []*discordgo.MessageEmbed{}
	embeds = append(embeds, &message.MessageEmbed)
	return ctx.CmdReplyEmbeds(embeds)
}
func (ctx *CommandCtx) CmdReplyEmbeds(embeds []*discordgo.MessageEmbed) *discordgo.InteractionResponse {
	return &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: embeds,
		},
	}
}
func (ctx *CommandCtx) Reply(message string) *discordgo.InteractionResponse {
	return &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: message,
		},
	}
}
func (ctx *CommandCtx) Replyf(format string, a ...interface{}) *discordgo.InteractionResponse {
	return &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: fmt.Sprintf(format, a...),
		},
	}
}

func (ctx *CommandCtx) Errorf(format string, a ...interface{}) *discordgo.InteractionResponse {
	return ctx.Replyf(format, a...)
}
