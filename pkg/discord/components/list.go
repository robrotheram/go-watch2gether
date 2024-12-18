package components

import (
	"fmt"
	"math"
	"unicode"
	"w2g/pkg/controllers"
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

func paginate(pageNum int, pageSize int, sliceLength int) (int, int) {
	start := pageNum * pageSize
	if start > sliceLength {
		start = sliceLength
	}
	end := start + pageSize
	if end > sliceLength {
		end = sliceLength
	}

	if start < 0 {
		start = 0
	}
	if end > sliceLength {
		end = sliceLength - 1
	}
	return start, end
}

func maxPages(queue []*media.Media) int {
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

func emptyQueue() *discordgo.InteractionResponseData {
	return &discordgo.InteractionResponseData{
		Content: "Queue is empty",
		Flags:   discordgo.MessageFlagsEphemeral,
	}
}

func QueueCompontent(state *controllers.PlayerState, pageNum int) *discordgo.InteractionResponseData {
	if len(state.Queue) == 0 {
		return emptyQueue()
	}

	start, end := paginate(pageNum, pageSize, len(state.Queue))
	pagedSlice := state.Queue[start:end]
	if len(pagedSlice) == 0 {
		pageNum = pageNum - 1
		start, end := paginate(pageNum, pageSize, len(state.Queue))
		pagedSlice = state.Queue[start:end]
	}

	embed := EmbedBuilder("Currently in the Queue")
	queStr := ""
	for i, video := range pagedSlice {
		pos := pageNum*pageSize + i + 1
		queStr = queStr + fmt.Sprintf("`%d.` [%s](%s) \n", pos, truncate(video.Title, 40), video.Url)
	}
	totalPages := float64(len(state.Queue)) / float64(pageSize)

	embed.AddField(discordgo.MessageEmbedField{
		Name:  fmt.Sprintf("Page %d of %d", pageNum+1, int(math.Ceil(totalPages))),
		Value: queStr,
	})
	embed.Description = fmt.Sprintf("%d tracks in total in the queue", len(state.Queue))

	components := []discordgo.MessageComponent{
		discordgo.ActionsRow{
			Components: []discordgo.MessageComponent{
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
			},
		},
	}

	currentlyPlayingEmbed := EmbedBuilder("Nothing is currently playing")
	if state.Current != nil {
		currentlyPlayingEmbed = MediaEmbed(*state.Current, "Currently Playing:")
	}

	return &discordgo.InteractionResponseData{
		Flags:      discordgo.MessageFlagsEphemeral,
		Components: components,
		Embeds: []*discordgo.MessageEmbed{
			&currentlyPlayingEmbed.MessageEmbed,
			&embed.MessageEmbed,
		},
	}
}
