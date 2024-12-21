package components

import (
	"fmt"
	"w2g/pkg/controllers"

	"github.com/bwmarrin/discordgo"
)

var (
	QueueBtnFirst    = "QUUUE_BTN_FIRST"
	QueueBtnLast     = "QUUUE_BTN_Last"
	QueueBtnNext     = "QUUUE_BTN_NEXT"
	QueueBtnPrevious = "QUUUE_BTN_PRE"
)

func init() {
	register(
		Handler{
			Name: QueueBtnFirst,
			Function: func(ctx HandlerCtx) *discordgo.InteractionResponse {
				ctx.UserSession.Page = 0
				return ctx.UpdateMessage(QueueCompontent(ctx.Controller.State(), ctx.UserSession.Page))
			},
		},
	)
	register(
		Handler{
			Name: QueueBtnLast,
			Function: func(ctx HandlerCtx) *discordgo.InteractionResponse {
				ctx.UserSession.Page = maxPages(ctx.Controller.State().Queue)
				return ctx.UpdateMessage(QueueCompontent(ctx.Controller.State(), ctx.UserSession.Page))
			},
		},
	)
	register(
		Handler{
			Name: QueueBtnNext,
			Function: func(ctx HandlerCtx) *discordgo.InteractionResponse {
				ctx.UserSession.Page = ctx.UserSession.Page + 1
				if ctx.UserSession.Page > maxPages(ctx.Controller.State().Queue) {
					ctx.UserSession.Page = maxPages(ctx.Controller.State().Queue)
				}
				return ctx.UpdateMessage(QueueCompontent(ctx.Controller.State(), ctx.UserSession.Page))
			},
		},
	)
	register(
		Handler{
			Name: QueueBtnPrevious,
			Function: func(ctx HandlerCtx) *discordgo.InteractionResponse {
				ctx.UserSession.Page = ctx.UserSession.Page - 1
				if ctx.UserSession.Page < 0 {
					ctx.UserSession.Page = 0
				}
				return ctx.UpdateMessage(QueueCompontent(ctx.Controller.State(), ctx.UserSession.Page))
			},
		},
	)
}

func QueueCompontent(state *controllers.PlayerState, pageNum int) *discordgo.InteractionResponseData {
	if len(state.Queue) == 0 {
		return emptylist("Queue is empty")
	}

	buttons := []discordgo.MessageComponent{
		discordgo.Button{
			CustomID: QueueBtnFirst,
			Emoji: &discordgo.ComponentEmoji{
				Name: "⏮️",
			},
			Style: discordgo.PrimaryButton,
		},
		discordgo.Button{
			CustomID: QueueBtnPrevious,
			Emoji: &discordgo.ComponentEmoji{
				Name: "◀️",
			},
			Style: discordgo.PrimaryButton,
		},
		discordgo.Button{
			CustomID: QueueBtnNext,
			Emoji: &discordgo.ComponentEmoji{
				Name: "▶️",
			},
			Style: discordgo.PrimaryButton,
		},
		discordgo.Button{
			CustomID: QueueBtnLast,
			Emoji: &discordgo.ComponentEmoji{
				Name: "⏭️",
			},
			Style: discordgo.PrimaryButton,
		},
	}

	cmp := ListCompontent(
		state.Queue,
		pageNum,
		buttons,
		"Currently in the Queue",
		fmt.Sprintf("%d tracks in total in the queue", len(state.Queue)),
	)

	currentlyPlayingEmbed := EmbedBuilder("Nothing is currently playing")
	if state.Current != nil {
		currentlyPlayingEmbed = MediaEmbed(*state.Current, "Currently Playing:")
	}
	cmp.Embeds = append([]*discordgo.MessageEmbed{&currentlyPlayingEmbed.MessageEmbed}, cmp.Embeds...)
	return cmp
}
