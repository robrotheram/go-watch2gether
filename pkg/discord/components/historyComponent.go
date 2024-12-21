package components

import (
	"fmt"
	"w2g/pkg/media"

	"github.com/bwmarrin/discordgo"
)

var (
	HistoryBtnFirst    = "HISTORY_BTN_FIRST"
	HistoryBtnLast     = "HISTORY_BTN_Last"
	HistoryBtnNext     = "HISTORY_BTN_NEXT"
	HistoryBtnPrevious = "HISTORY_BTN_PRE"
)

func init() {
	register(
		Handler{
			Name: HistoryBtnFirst,
			Function: func(ctx HandlerCtx) *discordgo.InteractionResponse {
				ctx.UserSession.Page = 0
				history, _ := ctx.Controller.History()
				return ctx.UpdateMessage(HistoryCompontent(history, ctx.UserSession.Page))
			},
		},
	)
	register(
		Handler{
			Name: HistoryBtnLast,
			Function: func(ctx HandlerCtx) *discordgo.InteractionResponse {
				history, _ := ctx.Controller.History()
				ctx.UserSession.Page = maxPages(history)
				return ctx.UpdateMessage(HistoryCompontent(history, ctx.UserSession.Page))
			},
		},
	)
	register(
		Handler{
			Name: HistoryBtnNext,
			Function: func(ctx HandlerCtx) *discordgo.InteractionResponse {
				history, _ := ctx.Controller.History()
				ctx.UserSession.Page = ctx.UserSession.Page + 1
				if ctx.UserSession.Page > maxPages(history) {
					ctx.UserSession.Page = maxPages(history)
				}
				return ctx.UpdateMessage(HistoryCompontent(history, ctx.UserSession.Page))
			},
		},
	)
	register(
		Handler{
			Name: QueueBtnPrevious,
			Function: func(ctx HandlerCtx) *discordgo.InteractionResponse {
				history, _ := ctx.Controller.History()
				ctx.UserSession.Page = ctx.UserSession.Page - 1
				if ctx.UserSession.Page < 0 {
					ctx.UserSession.Page = 0
				}
				return ctx.UpdateMessage(HistoryCompontent(history, ctx.UserSession.Page))
			},
		},
	)
}

func HistoryCompontent(history []*media.Media, pageNum int) *discordgo.InteractionResponseData {

	if len(history) == 0 {
		return emptylist("Nothing in the hisory")
	}

	buttons := []discordgo.MessageComponent{
		discordgo.Button{
			CustomID: HistoryBtnFirst,
			Emoji: &discordgo.ComponentEmoji{
				Name: "⏮️",
			},
			Style: discordgo.PrimaryButton,
		},
		discordgo.Button{
			CustomID: HistoryBtnPrevious,
			Emoji: &discordgo.ComponentEmoji{
				Name: "◀️",
			},
			Style: discordgo.PrimaryButton,
		},
		discordgo.Button{
			CustomID: HistoryBtnNext,
			Emoji: &discordgo.ComponentEmoji{
				Name: "▶️",
			},
			Style: discordgo.PrimaryButton,
		},
		discordgo.Button{
			CustomID: HistoryBtnLast,
			Emoji: &discordgo.ComponentEmoji{
				Name: "⏭️",
			},
			Style: discordgo.PrimaryButton,
		},
	}

	return ListCompontent(
		history,
		pageNum,
		buttons,
		"Media Timeline ",
		fmt.Sprintf("%d tracks in the timeline", len(history)),
	)
}
