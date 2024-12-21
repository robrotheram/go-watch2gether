package components

import (
	"fmt"
	"math"
	"unicode"
	"w2g/pkg/media"

	"github.com/bwmarrin/discordgo"
)

type ButtonID string

const pageSize = 10

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

func emptylist(msg string) *discordgo.InteractionResponseData {
	return &discordgo.InteractionResponseData{
		Content: msg,
		Flags:   discordgo.MessageFlagsEphemeral,
	}
}

func ListCompontent(list []*media.Media, pageNum int, buttons []discordgo.MessageComponent, title string, description string) *discordgo.InteractionResponseData {
	start, end := paginate(pageNum, pageSize, len(list))
	pagedSlice := list[start:end]
	if len(pagedSlice) == 0 {
		pageNum = pageNum - 1
		start, end := paginate(pageNum, pageSize, len(list))
		pagedSlice = list[start:end]
	}

	embed := EmbedBuilder(title)
	queStr := ""
	for i, video := range pagedSlice {
		pos := pageNum*pageSize + i + 1
		queStr = queStr + fmt.Sprintf("`%d.` [%s](%s) \n", pos, truncate(video.Title, 40), video.Url)
	}
	totalPages := float64(len(list)) / float64(pageSize)

	embed.AddField(discordgo.MessageEmbedField{
		Name:  fmt.Sprintf("Page %d of %d", pageNum+1, int(math.Ceil(totalPages))),
		Value: queStr,
	})
	embed.Description = description

	components := []discordgo.MessageComponent{
		discordgo.ActionsRow{
			Components: buttons,
		},
	}

	return &discordgo.InteractionResponseData{
		Flags:      discordgo.MessageFlagsEphemeral,
		Components: components,
		Embeds: []*discordgo.MessageEmbed{
			&embed.MessageEmbed,
		},
	}
}
