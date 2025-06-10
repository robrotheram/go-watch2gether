package commands

import (
	"w2g/pkg/utils"

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
			Name: "seek",
			ApplicationCommand: []discordgo.ApplicationCommand{
				{
					Description: "Set the position of the track to the given time. ",
					Options: []*discordgo.ApplicationCommandOption{
						{
							Type:        discordgo.ApplicationCommandOptionString,
							Name:        "time",
							Description: "Position to fast forward e.g 30 for seconds ",
							Required:    true,
						},
					},
				},
			},
			Function: seekCMD,
		},
		Command{
			Name: "restart",
			ApplicationCommand: []discordgo.ApplicationCommand{
				{
					Description: "Restart the currently playing track.",
				},
			},
			Function: restartCmd,
		},
		Command{
			Name: "loop",
			ApplicationCommand: []discordgo.ApplicationCommand{
				{
					Description: "Enable or disabling Looping.",
				},
			},
			Function: loopCMD,
		},
	)
}

func getPostionOption(ctx CommandCtx) string {
	if len(ctx.Args) < 2 {
		return ""
	}
	return ctx.Args[1]
}

func playCmd(ctx CommandCtx) *discordgo.InteractionResponse {
	if !ctx.Controller.ContainsPlayer(ctx.Guild.ID) {
		if join(ctx) != nil {
			return ctx.Reply("User not connected to voice channel")
		}
	}
	go ctx.Controller.Start(ctx.Member.User.Username)
	return ctx.Reply(":play_pause: Now Playing :thumbsup:")
}

func seekCMD(ctx CommandCtx) *discordgo.InteractionResponse {
	if !ctx.Controller.ContainsPlayer(ctx.Guild.ID) {
		if join(ctx) != nil {
			return ctx.Reply("User not connected to voice channel")
		}
	}
	seekTime, err := utils.ParseTime(ctx.Args[0])
	if err != nil {
		return ctx.Reply("Invalid time format")
	}
	go ctx.Controller.Seek(seekTime, ctx.Member.User.Username)
	return ctx.Replyf(":fast_forward: Seeking to %f seconds into the track :thumbsup:", seekTime.Seconds())
}

func restartCmd(ctx CommandCtx) *discordgo.InteractionResponse {
	if !ctx.Controller.ContainsPlayer(ctx.Guild.ID) {
		if join(ctx) != nil {
			return ctx.Reply("User not connected to voice channel")
		}
	}
	go ctx.Controller.Seek(0, ctx.Member.User.Username)
	return ctx.Reply(":leftwards_arrow_with_hook: Restarting Track")
}

func skipCmd(ctx CommandCtx) *discordgo.InteractionResponse {
	go ctx.Controller.Skip(ctx.Member.User.Username)
	return ctx.Reply(":fast_forward: Now Skipping :thumbsup:")
}

func stopCmd(ctx CommandCtx) *discordgo.InteractionResponse {
	go ctx.Controller.Stop(ctx.Member.User.Username)
	return ctx.Reply(":stop_button: Stopping track")
}

func pauseCmd(ctx CommandCtx) *discordgo.InteractionResponse {
	go ctx.Controller.Pause(ctx.Member.User.Username)
	return ctx.Reply("Pause")
}

func loopCMD(ctx CommandCtx) *discordgo.InteractionResponse {
	ctx.Controller.Loop(ctx.Member.User.Username)
	if ctx.Controller.State().Loop {
		return ctx.Reply(":arrows_counterclockwise: Loop On")
	}
	return ctx.Reply(":arrows_counterclockwise: Loop Off")
}
