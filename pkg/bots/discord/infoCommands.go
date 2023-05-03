package discordbot

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

func init() {
	Commands.Register(
		CMD{
			ApplicationCommand: discordgo.ApplicationCommand{
				Name:        "watch",
				Description: "Get Url Watch2gether Room, where the video will be in sync with discord",
			},
			Function: LinkCMD,
		},

		CMD{
			ApplicationCommand: discordgo.ApplicationCommand{
				Name:        "nowplaying",
				Description: "Shows what song is currently playing",
			},
			Function: nowPlayingCmd,
		},

		CMD{
			ApplicationCommand: discordgo.ApplicationCommand{
				Name:        "queue",
				Description: "Shows what is in the current queue",
			},
			Function: queueCMD,
		},
		// CMD{
		// 	ApplicationCommand: discordgo.ApplicationCommand{
		// 		Name:        "version",
		// 		Description: "Watch2Gether Version",
		// 	},
		// 	Function: VersionCMD,
		// },
	)
}

func LinkCMD(ctx CommandCtx) *discordgo.InteractionResponse {
	msg := EmbedBuilder("Watch2Gether")
	msg.URL = fmt.Sprintf("%s/app/room/%s", ctx.BaseURL, ctx.Guild.ID)
	msg.Type = discordgo.EmbedTypeArticle
	msg.Description = fmt.Sprintf("%s/app/room/%s", ctx.BaseURL, ctx.Guild.ID)
	return ctx.CmdReplyEmbed(msg)
}

func nowPlayingCmd(ctx CommandCtx) *discordgo.InteractionResponse {

	r, err := ctx.GetChannel(ctx.Guild.ID)
	if err != nil {
		return ctx.Errorf("Room %s not active", ctx.Guild.ID)
	}
	var video = r.GetCurrentVideo()

	message := EmbedBuilder("Currently Playing")
	message.Thumbnail = &discordgo.MessageEmbedThumbnail{
		URL: video.Thumbnail,
	}
	message.Description = video.Title

	message.AddField(discordgo.MessageEmbedField{
		Name:   "Channel",
		Value:  video.ChannelName,
		Inline: true,
	})
	message.AddField(discordgo.MessageEmbedField{
		Name:   "Current Progress",
		Value:  r.Duration().String(),
		Inline: true,
	})

	message.AddField(discordgo.MessageEmbedField{
		Name:   "Song Duration",
		Value:  video.Duration.String(),
		Inline: true,
	})
	return ctx.CmdReplyEmbed(message)
}

func queueCMD(ctx CommandCtx) *discordgo.InteractionResponse {
	r, err := ctx.GetChannel(ctx.Guild.ID)
	if err != nil {
		return ctx.Errorf("Room %s not active", ctx.Guild.ID)
	}
	msg := EmbedBuilder("Watch2Gether Queue")
	msg.Thumbnail = nil

	video := r.GetCurrentVideo()
	queue := r.GetQueue()

	if video.Url != "" {
		msg.AddField(discordgo.MessageEmbedField{
			Name: "Now Playing:",
			Value: fmt.Sprintf(
				"[%s](%s) | `%s Requested by: %s`",
				video.Title,
				video.Url,
				video.Duration,
				video.User),
		})
	}

	queStr := ""
	for i, video := range queue {
		if i >= 5 {
			break
		}
		queStr = queStr + fmt.Sprintf("`%d.` [%s](%s) | `%s Requested by: %s` \n\n",
			i+1,
			video.Title,
			video.Url,
			video.Duration,
			video.User)
	}

	if len(queStr) == 0 {
		queStr = "There is nothing the the queue"
	}
	msg.AddField(discordgo.MessageEmbedField{
		Name:  "Up Next:",
		Value: queStr,
	})

	msg.Description = fmt.Sprintf("%d tracks in total in the queue", len(queue))
	return ctx.CmdReplyEmbed(msg)
}

// func VersionCMD(ctx CommandCtx) *discordgo.InteractionResponse {
// 	return ctx.Reply(fmt.Sprintf("Version: %s", datastore.VERSION))
// }
