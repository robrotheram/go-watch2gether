package components

import (
	"w2g/pkg/media"

	"github.com/bwmarrin/discordgo"
)

func MediaEmbed(media media.Media, title string) *EmbededMessage {
	embed := EmbedBuilder(title)
	embed.Thumbnail = &discordgo.MessageEmbedThumbnail{
		URL: media.Thumbnail,
	}
	embed.Description = media.Title

	embed.AddField(discordgo.MessageEmbedField{
		Name:   "Channel",
		Value:  media.ChannelName,
		Inline: true,
	})
	embed.AddField(discordgo.MessageEmbedField{
		Name:   "Song Duration",
		Value:  media.Progress.Duration.String(),
		Inline: true,
	})
	return embed
}
