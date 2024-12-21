package commands

import (
	"time"
	"w2g/pkg/discord/components"
	"w2g/pkg/utils"

	"github.com/bwmarrin/discordgo"
)

func init() {
	register(
		Command{
			Name: "controls",
			ApplicationCommand: []discordgo.ApplicationCommand{
				{
					Type: discordgo.UserApplicationCommand,
				},
				{
					Description: "show player controls",
					Type:        discordgo.ChatApplicationCommand,
				},
			},
			Function: controlscmd,
		},
		Command{
			Name: "now",
			ApplicationCommand: []discordgo.ApplicationCommand{
				{
					Type: discordgo.UserApplicationCommand,
				},
				{
					Description: "show what is currently playing",
					Type:        discordgo.ChatApplicationCommand,
				},
			},
			Function: controlscmd,
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
		},
		Command{
			Name: "version",
			ApplicationCommand: []discordgo.ApplicationCommand{
				{
					Description: "return the version of this bot",
					Type:        discordgo.ChatApplicationCommand,
				},
			},
			Function: versioncmd,
		},
		Command{
			Name: "history",
			ApplicationCommand: []discordgo.ApplicationCommand{
				{
					Description: "what was previously played in this channel",
					Type:        discordgo.ChatApplicationCommand,
				},
			},
			Function: historycmd,
		},
	)
}

func listcmd(ctx CommandCtx) *discordgo.InteractionResponse {
	return ctx.CmdReplyData(components.QueueCompontent(ctx.Controller.State(), 0))
}

func controlscmd(ctx CommandCtx) *discordgo.InteractionResponse {
	return ctx.CmdReplyData(components.ControlCompontent(ctx.Controller.State()))
}

func versioncmd(ctx CommandCtx) *discordgo.InteractionResponse {
	return ctx.Replyf("Version: %s commit: %s \n time: %s", utils.Version, utils.Revision, utils.LastCommit.Format(time.RFC822))
}

func historycmd(ctx CommandCtx) *discordgo.InteractionResponse {
	history, _ := ctx.Controller.History()
	return ctx.CmdReplyData(components.HistoryCompontent(history, 0))
}
