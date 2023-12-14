package components

import (
	"fmt"
	"unicode"
	"w2g/pkg/media"

	"github.com/bwmarrin/discordgo"
)

type ButtonID string

var (
	QueueBtnFirst    = "QUUUE_BTN_FIRST"
	QueueBtnLast     = "QUUUE_BTN_Last"
	QueueBtnNext     = "QUUUE_BTN_NEXT"
	QueueBtnPrevious = "QUUUE_BTN_PRE"
)

var pageSize = 10

func init() {
	register(
		Handler{
			Name: QueueBtnFirst,
			Function: func(ctx HandlerCtx) *discordgo.InteractionResponse {
				ctx.UserSession.Page = 0
				return ctx.UpdateMessage(QueueCompontent(ctx.Controller.State().Queue, ctx.UserSession.Page))
			},
		},
	)
	register(
		Handler{
			Name: QueueBtnLast,
			Function: func(ctx HandlerCtx) *discordgo.InteractionResponse {
				ctx.UserSession.Page = maxPages(ctx.Controller.State().Queue)
				return ctx.UpdateMessage(QueueCompontent(ctx.Controller.State().Queue, ctx.UserSession.Page))
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
				return ctx.UpdateMessage(QueueCompontent(ctx.Controller.State().Queue, ctx.UserSession.Page))
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
				return ctx.UpdateMessage(QueueCompontent(ctx.Controller.State().Queue, ctx.UserSession.Page))
			},
		},
	)
}

func paginate(pageNum int, pageSize int, sliceLength int) (int, int) {
	start := pageNum * pageSize
	if start > sliceLength {
		start = sliceLength
	}
	end := start + pageSize
	if end > sliceLength {
		end = sliceLength
	}
	return start, end
}

func maxPages(queue []media.Media) int {
	return len(queue) / 10
}

func truncate(text string, maxLen int) string {
	lastSpaceIx := maxLen
	len := 0
	for i, r := range text {
		if unicode.IsSpace(r) {
			lastSpaceIx = i
		}
		len++
		if len > maxLen {
			return text[:lastSpaceIx] + "..."
		}
	}
	return text
}

func QueueCompontent(queue []media.Media, pageNum int) *discordgo.InteractionResponseData {
	start, end := paginate(pageNum, pageSize, len(queue))
	pagedSlice := queue[start:end]

	embed := EmbedBuilder("Currently in the Queue")

	queStr := ""
	for i, video := range pagedSlice {
		pos := pageNum*pageSize + i + 1
		queStr = queStr + fmt.Sprintf("`%d.` [%s](%s) \n", pos, truncate(video.Title, 40), video.Url)
	}
	if len(queStr) == 0 {
		queStr = "There is nothing the the queue"
	}

	embed.AddField(discordgo.MessageEmbedField{
		Name:  "Up Next:",
		Value: queStr,
	})
	embed.Description = fmt.Sprintf("%d tracks in total in the queue", len(queue))

	components := []discordgo.MessageComponent{
		discordgo.ActionsRow{
			Components: []discordgo.MessageComponent{
				discordgo.Button{
					CustomID: QueueBtnFirst,
					Emoji: discordgo.ComponentEmoji{
						Name: "⏮️",
					},
					Style: discordgo.PrimaryButton,
				},
				discordgo.Button{
					CustomID: QueueBtnPrevious,
					Emoji: discordgo.ComponentEmoji{
						Name: "◀️",
					},
					Style: discordgo.PrimaryButton,
				},
				discordgo.Button{
					CustomID: QueueBtnNext,
					Emoji: discordgo.ComponentEmoji{
						Name: "▶️",
					},
					Style: discordgo.PrimaryButton,
				},
				discordgo.Button{
					CustomID: QueueBtnLast,
					Emoji: discordgo.ComponentEmoji{
						Name: "⏭️",
					},
					Style: discordgo.PrimaryButton,
				},
			},
		},
	}

	return &discordgo.InteractionResponseData{
		Content:    "Currently Playing:",
		Flags:      discordgo.MessageFlagsEphemeral,
		Components: components,
		Embeds: []*discordgo.MessageEmbed{
			&embed.MessageEmbed,
		},
	}
}
