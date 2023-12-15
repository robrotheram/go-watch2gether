package commands

import (
	"w2g/pkg/discord/components"

	"github.com/bwmarrin/discordgo"
)

func init() {
	register(
		Command{
			ApplicationCommand: discordgo.ApplicationCommand{
				Name: "nowplaying",
				Type: discordgo.UserApplicationCommand,
			},
			Function: nowplayingCmd,
		},
		Command{
			ApplicationCommand: discordgo.ApplicationCommand{
				Name:        "nowplaying",
				Description: "show what is currently playing",
				Type:        discordgo.ChatApplicationCommand,
			},
			Function: nowplayingCmd,
		},
		Command{
			ApplicationCommand: discordgo.ApplicationCommand{
				Name:        "list",
				Description: "list what tracks are currently in the queue",
				Type:        discordgo.ChatApplicationCommand,
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
