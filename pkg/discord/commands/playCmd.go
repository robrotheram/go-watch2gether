package commands

import (
	"net/url"

	"github.com/bwmarrin/discordgo"
)

func init() {
	register(
		Command{
			ApplicationCommand: discordgo.ApplicationCommand{
				Name:        "play",
				Description: "plays what is currently in the queue",
				Type:        discordgo.ChatApplicationCommand,
			},
			Function: playCmd,
		},
		Command{
			ApplicationCommand: discordgo.ApplicationCommand{
				Name:        "skip",
				Description: "Skip the current track",
				Type:        discordgo.ChatApplicationCommand,
			},
			Function: skipCmd,
		},
		Command{
			ApplicationCommand: discordgo.ApplicationCommand{
				Name:        "stop",
				Description: "Stop playing",
				Type:        discordgo.ChatApplicationCommand,
			},
			Function: stopCmd,
		},
		Command{
			ApplicationCommand: discordgo.ApplicationCommand{
				Name:        "pause",
				Description: "Pause the current track",
				Type:        discordgo.ChatApplicationCommand,
			},
			Function: pauseCmd,
		},
		Command{
			ApplicationCommand: discordgo.ApplicationCommand{
				Name:        "add",
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
			Function: addCmd,
		})
}

func addCmd(ctx CommandCtx) *discordgo.InteractionResponse {
	_, err := url.ParseRequestURI(ctx.Args[0])
	if err != nil {
		return ctx.Errorf("%s Is not a valid URL", ctx.Args[0])
	}
	err = ctx.Controller.Add(ctx.Args[0], "")
	if err != nil {
		return ctx.Errorf("error: %v", err)
	}
	return ctx.Reply(":notes: added traks to the queue :thumbsup:")
}

func playCmd(ctx CommandCtx) *discordgo.InteractionResponse {
	if !ctx.Controller.IsActive() {
		if join(ctx) != nil {
			return ctx.Reply("User not connected to voice channel")
		}
	}
	ctx.Controller.Start()
	return ctx.Reply(":play_pause: Now Playing :thumbsup:")
}

func skipCmd(ctx CommandCtx) *discordgo.InteractionResponse {
	ctx.Controller.Skip()
	return ctx.Reply(":fast_forward: Now Skipping :thumbsup:")
}

func stopCmd(ctx CommandCtx) *discordgo.InteractionResponse {
	ctx.Controller.Stop()
	return ctx.Reply(":stop_button: Stopping track")
}

func pauseCmd(ctx CommandCtx) *discordgo.InteractionResponse {
	ctx.Controller.Pause()
	return ctx.Reply("Pause")
}
