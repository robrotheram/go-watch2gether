package commands

import (
	"w2g/pkg/discord/components"

	"github.com/bwmarrin/discordgo"
)

func init() {
	register(
		Command{
			Name: "nowplaying",
			ApplicationCommand: []discordgo.ApplicationCommand{
				{

					Type: discordgo.UserApplicationCommand,
				},
				{
					Description: "show what is currently playing",
					Type:        discordgo.ChatApplicationCommand,
				},
			},
			Function: nowplayingCmd,
		},
		Command{
			Name: "list",
			ApplicationCommand: []discordgo.ApplicationCommand{
				{
					Description: "list what tracks are currently in the queue",
					Type:        discordgo.ChatApplicationCommand,
				},
			},
			Function: listcmd,
		})
}

func nowplayingCmd(ctx CommandCtx) *discordgo.InteractionResponse {
	current := ctx.Controller.State().Current
	return ctx.CmdReplyEmbed(components.MediaEmbed(current, "Now Playing"))
}

func listcmd(ctx CommandCtx) *discordgo.InteractionResponse {
	return ctx.CmdReplyData(components.QueueCompontent(ctx.Controller.State().Queue, 0))
}
