package components

import "github.com/bwmarrin/discordgo"

type EmbededMessage struct {
	discordgo.MessageEmbed
}

func EmbedBuilder(title string) *EmbededMessage {
	embed := EmbededMessage{}
	embed.Title = title
	embed.Thumbnail = &discordgo.MessageEmbedThumbnail{
		URL:    "https://watch2gether.exceptionerror.io/static/media/logo.14caa8da.jpg",
		Width:  10,
		Height: 10,
	}
	embed.Color = 0x4286f4
	return &embed
}

func (embed *EmbededMessage) AddField(field discordgo.MessageEmbedField) {
	if len(field.Value) > 1000 {
		field.Value = field.Value[:1000]
	}
	embed.Fields = append(embed.Fields, &field)
}
