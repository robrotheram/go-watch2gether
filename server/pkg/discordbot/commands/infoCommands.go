package commands

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

func init() {
	Register(
		CMD{
			Command:     "watch",
			Description: "Get Url Watch2gether Room, where the video will be in sync with discord",
			Aliases:     []string{"link"},
			Function:    LinkCMD,
		},

		CMD{
			Command:     "nowplaying",
			Description: "Shows what song is currently playing",
			Function:    nowPlayingCmd,
		},

		CMD{
			Command:     "queue",
			Description: "Other Usage: !queue <page>: Shows the specified page number.",
			Aliases:     []string{"q"},
			Function:    queueCMD,
		},
	)
}

func LinkCMD(ctx CommandCtx) error {
	msg := EmbedBuilder("Watch2Gether")
	msg.URL = fmt.Sprintf("%s/app/room/%s", ctx.BaseURL, ctx.Guild.ID)
	msg.Type = discordgo.EmbedTypeArticle
	msg.Description = fmt.Sprintf("%s/app/room/%s", ctx.BaseURL, ctx.Guild.ID)
	return ctx.ReplyEmbed(msg)
}

func nowPlayingCmd(ctx CommandCtx) error {

	meta, err := ctx.GetMeta()
	if err != nil {
		return fmt.Errorf("Room %s not active", ctx.Guild.ID)
	}
	video := meta.CurrentVideo

	message := EmbedBuilder("Currently Playing")
	message.Thumbnail = &discordgo.MessageEmbedThumbnail{
		URL: video.Thumbnail,
	}
	message.Description = video.Title

	message.AddField(discordgo.MessageEmbedField{
		Name:   "Channel",
		Value:  video.Channel,
		Inline: true,
	})

	message.AddField(discordgo.MessageEmbedField{
		Name:   "Song Duration",
		Value:  video.Duration.String(),
		Inline: true,
	})

	return ctx.ReplyEmbed(message)

}

func queueCMD(ctx CommandCtx) error {
	meta, err := ctx.GetMeta()
	if err != nil {
		return fmt.Errorf("Room %s not active", ctx.Guild.ID)
	}

	msg := EmbedBuilder("Watch2Gether Queue")
	msg.Thumbnail = nil
	msg.AddField(discordgo.MessageEmbedField{
		Name: "Now Playing",
		Value: fmt.Sprintf(
			"[%s](%s) | `%s Requested by: %s`",
			meta.CurrentVideo.Title,
			meta.CurrentVideo.Url,
			meta.CurrentVideo.Duration,
			meta.CurrentVideo.User),
	})

	queStr := ""
	for i, vidoe := range meta.Queue {
		if i >= 5 {
			break
		}
		queStr = queStr + fmt.Sprintf("`%d.` [%s](%s) | `%s Requested by: %s` \n\n",
			i+1,
			vidoe.Title,
			vidoe.Url,
			vidoe.Duration,
			vidoe.User)
	}

	msg.AddField(discordgo.MessageEmbedField{
		Name:  "Up Next:",
		Value: queStr,
	})

	msg.Description = fmt.Sprintf("%d tracks in total in the queue", len(meta.Queue))

	return ctx.ReplyEmbed(msg)

}
