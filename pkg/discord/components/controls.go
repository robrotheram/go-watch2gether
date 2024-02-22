package components

import (
	"w2g/pkg/controllers"

	"github.com/bwmarrin/discordgo"
)

var (
	PlayBtn  = "PLAY_BTN"
	PauseBtn = "PAUSE_BTN"
	SkipBtn  = "SKIP_BTN"
	StopBtn  = "STOP_BTN"
)

func init() {
	register(
		Handler{
			Name: PlayBtn,
			Function: func(ctx HandlerCtx) *discordgo.InteractionResponse {
				ctx.Controller.Start(ctx.User.Username)
				return ctx.UpdateMessage(ControlCompontent(ctx.Controller.State()))
			},
		},
	)
	register(
		Handler{
			Name: PauseBtn,
			Function: func(ctx HandlerCtx) *discordgo.InteractionResponse {
				ctx.Controller.Pause(ctx.User.Username)
				return ctx.UpdateMessage(ControlCompontent(ctx.Controller.State()))
			},
		},
	)
	register(
		Handler{
			Name: SkipBtn,
			Function: func(ctx HandlerCtx) *discordgo.InteractionResponse {
				ctx.Controller.Skip(ctx.User.Username)
				state := ctx.Controller.State()
				state.Next()
				return ctx.UpdateMessage(ControlCompontent(state))
			},
		},
	)
	register(
		Handler{
			Name: StopBtn,
			Function: func(ctx HandlerCtx) *discordgo.InteractionResponse {
				ctx.Controller.Stop(ctx.User.Username)
				return ctx.UpdateMessage(ControlCompontent(ctx.Controller.State()))
			},
		},
	)
}

func ControlCompontent(state *controllers.PlayerState) *discordgo.InteractionResponseData {
	var actionButton discordgo.Button

	if state.State == controllers.PLAY && state.Current.ID != "" {
		actionButton = discordgo.Button{
			CustomID: PauseBtn,
			Emoji: discordgo.ComponentEmoji{
				Name: "⏸️",
			},
			Style: discordgo.PrimaryButton,
		}
	} else {
		actionButton = discordgo.Button{
			CustomID: PlayBtn,
			Emoji: discordgo.ComponentEmoji{
				Name: "▶️",
			},
			Style: discordgo.PrimaryButton,
		}
	}

	components := []discordgo.MessageComponent{
		discordgo.ActionsRow{
			Components: []discordgo.MessageComponent{
				discordgo.Button{
					CustomID: StopBtn,
					Emoji: discordgo.ComponentEmoji{
						Name: "⏹️",
					},
					Style: discordgo.PrimaryButton,
				},
				actionButton,
				discordgo.Button{
					CustomID: SkipBtn,
					Emoji: discordgo.ComponentEmoji{
						Name: "⏭️",
					},
					Style: discordgo.PrimaryButton,
				},
			},
		},
	}

	currentlyPlayingEmbed := EmbedBuilder("Nothing is currently playing")
	if len(state.Current.ID) > 0 {
		currentlyPlayingEmbed = MediaEmbed(state.Current, "Currently Playing:")
	}

	return &discordgo.InteractionResponseData{
		Content:    "Currently Playing:",
		Flags:      discordgo.MessageFlagsEphemeral,
		Components: components,
		Embeds: []*discordgo.MessageEmbed{
			&currentlyPlayingEmbed.MessageEmbed,
		},
	}
}
