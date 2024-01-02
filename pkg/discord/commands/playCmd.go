package commands

import (
	"net/url"
	"w2g/pkg/discord/players"

	"github.com/bwmarrin/discordgo"
)

func init() {
	register(
		Command{
			Name: "play",
			ApplicationCommand: []discordgo.ApplicationCommand{
				{
					Description: "plays what is currently in the queue",
					Type:        discordgo.ChatApplicationCommand,
				},
			},
			Function: playCmd,
		},
		Command{
			Name: "skip",
			ApplicationCommand: []discordgo.ApplicationCommand{
				{
					Description: "Skip the current track",
					Type:        discordgo.ChatApplicationCommand,
				},
			},
			Function: skipCmd,
		},
		Command{
			Name: "stop",
			ApplicationCommand: []discordgo.ApplicationCommand{
				{
					Description: "Stop playing",
					Type:        discordgo.ChatApplicationCommand,
				},
			},
			Function: stopCmd,
		},
		Command{
			Name: "pause",
			ApplicationCommand: []discordgo.ApplicationCommand{
				{
					Description: "Pause the current track",
					Type:        discordgo.ChatApplicationCommand,
				},
			},
			Function: pauseCmd,
		},
		Command{
			Name: "add",
			ApplicationCommand: []discordgo.ApplicationCommand{
				{
					Description: "Add new track to the queue",
					Options: []*discordgo.ApplicationCommandOption{
						{
							Type:        discordgo.ApplicationCommandOptionString,
							Name:        "media",
							Description: "Video/Audio URL e.g (https://www.youtube.com/watch?v=noneMROp_E8)",
							Required:    true,
						},
					},
					Type: discordgo.ChatApplicationCommand,
				},
			},
			Function: addCmd,
		})
}

func addCmd(ctx CommandCtx) *discordgo.InteractionResponse {
	_, err := url.ParseRequestURI(ctx.Args[0])
	if err != nil {
		return ctx.Errorf("%s Is not a valid URL", ctx.Args[0])
	}
	err = ctx.Controller.Add(ctx.Args[0], ctx.Member.User.Username)
	if err != nil {
		return ctx.Errorf("error: %v", err)
	}
	return ctx.Reply(":notes: added traks to the queue :thumbsup:")
}

func playCmd(ctx CommandCtx) *discordgo.InteractionResponse {
	if !ctx.Controller.ContainsPlayer(players.DISCORD) {
		if join(ctx) != nil {
			return ctx.Reply("User not connected to voice channel")
		}
	}
	ctx.Controller.Start(ctx.Member.User.Username)
	return ctx.Reply(":play_pause: Now Playing :thumbsup:")
}

func skipCmd(ctx CommandCtx) *discordgo.InteractionResponse {
	ctx.Controller.Skip(ctx.Member.User.Username)
	return ctx.Reply(":fast_forward: Now Skipping :thumbsup:")
}

func stopCmd(ctx CommandCtx) *discordgo.InteractionResponse {
	ctx.Controller.Stop(ctx.Member.User.Username)
	return ctx.Reply(":stop_button: Stopping track")
}

func pauseCmd(ctx CommandCtx) *discordgo.InteractionResponse {
	ctx.Controller.Pause(ctx.Member.User.Username)
	return ctx.Reply("Pause")
}
