package commands

import (
	"net/url"
	"w2g/pkg/controllers"

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
						{
							Type:        discordgo.ApplicationCommandOptionString,
							Name:        "top or bottom of the the queue",
							Description: "Where to add new media top of bottom of the queue (default: bottom)",
							Required:    false,
							Choices: []*discordgo.ApplicationCommandOptionChoice{
								{
									Name:  "top",
									Value: "TOP",
								},
								{
									Name:  "bottom",
									Value: "BOTTOM",
								},
							},
						},
					},
					Type: discordgo.ChatApplicationCommand,
				},
			},
			Function: addCmd,
		})
}

func getPostionOption(ctx CommandCtx) string {
	if len(ctx.Args) < 2 {
		return ""
	}
	return ctx.Args[1]
}

func addCmd(ctx CommandCtx) *discordgo.InteractionResponse {
	_, err := url.ParseRequestURI(ctx.Args[0])
	if err != nil {
		return ctx.Errorf("%s Is not a valid URL", ctx.Args[0])
	}
	isTop := getPostionOption(ctx) == "TOP"
	err = ctx.Controller.Add(ctx.Args[0], isTop, ctx.Member.User.Username)
	if err != nil {
		return ctx.Errorf("error: %v", err)
	}

	if !ctx.Controller.ContainsPlayer(ctx.Guild.ID) {
		join(ctx)
	}

	if ctx.Controller.State().State != controllers.PLAY {
		ctx.Controller.Start(ctx.Member.User.Username)
	}

	return ctx.Reply(":notes: added traks to the queue :thumbsup:")
}

func playCmd(ctx CommandCtx) *discordgo.InteractionResponse {
	if !ctx.Controller.ContainsPlayer(ctx.Guild.ID) {
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
