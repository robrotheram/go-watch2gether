package discordbot

import (
	"fmt"
	"net/url"
	"watch2gether/pkg/channels"
	"watch2gether/pkg/media"

	"github.com/bwmarrin/discordgo"
)

func init() {
	Commands.Register(
		CMD{
			ApplicationCommand: discordgo.ApplicationCommand{
				Name:        "play",
				Description: "play",
			},
			Function: playCmd,
		},
		CMD{
			ApplicationCommand: discordgo.ApplicationCommand{
				Name:        "add",
				Description: "add media(youtube, soundcloud, radio garden, mp3, mp4)",
				Options: []*discordgo.ApplicationCommandOption{
					{
						Type:        discordgo.ApplicationCommandOptionString,
						Name:        "media",
						Description: "Video/Audio URL e.g (https://www.youtube.com/watch?v=noneMROp_E8)",
						Required:    true,
					},
				},
			},
			Function: addCmd,
		},
		CMD{
			ApplicationCommand: discordgo.ApplicationCommand{
				Name:        "pause",
				Description: "Pauses the current playing track",
			},
			Aliases:  []string{"stop"},
			Function: pauseCMD,
		},
		CMD{
			ApplicationCommand: discordgo.ApplicationCommand{
				Name:        "skip",
				Description: "Skips the current song and plays the song you requested.",
			},
			Aliases:  []string{"fs", "skipped", "next"},
			Function: skipCMD,
		},
	)
}

func AddVideo(ctx CommandCtx, uri string) error {
	r, err := ctx.GetChannel(ctx.Guild.ID, channels.DISCORD)
	if err != nil {
		return fmt.Errorf("room %s not active", ctx.Guild.ID)
	}

	if media.MediaFactory.GetFactory(ctx.Args[0]) == nil {
		ctx.Reply(fmt.Sprintf("%s is a unsupported media type", ctx.Args[0]))
	}

	videos, err := media.NewVideo(uri, ctx.User.User.Username)
	if len(videos) == 0 || err != nil {
		return fmt.Errorf("no media found")
	}
	video := videos[0]
	queue, _ := r.GetQueue()
	message := EmbedBuilder(fmt.Sprintf("Added %d tracks to the Queue", len(videos)))
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
		Name:   "Song Duration",
		Value:  video.Duration.String(),
		Inline: true,
	})
	message.AddField(discordgo.MessageEmbedField{
		Name:   "Position in the Queue",
		Value:  fmt.Sprintf("%d", len(queue)+1),
		Inline: false,
	})

	r.Add(videos)
	return ctx.ReplyEmbed(message)
}

func addCmd(ctx CommandCtx) *discordgo.InteractionResponse {
	if len(ctx.Args) <= 0 {
		return ctx.Errorf("no media url found")
	}

	_, err := url.ParseRequestURI(ctx.Args[0])
	if err != nil {
		return ctx.Errorf("%s Is not a valid URL", ctx.Args[0])
	}

	_, err = ctx.GetChannel(ctx.Guild.ID, channels.DISCORD)
	if err != nil {
		return ctx.Errorf("Room %s not active", ctx.Guild.ID)
	}
	go AddVideo(ctx, ctx.Args[0])

	return ctx.Replyf("Proccessing media: %s, will be added to the queue shortly", ctx.Args[0])

}

func playCmd(ctx CommandCtx) *discordgo.InteractionResponse {
	r, err := ctx.GetChannel(ctx.Guild.ID, channels.DISCORD)
	if err != nil {
		return ctx.Errorf("Room %s not active", ctx.Guild.ID)
	}
	r.Play()
	return ctx.Reply(":play_pause: Now Playing :thumbsup:")
}

func pauseCMD(ctx CommandCtx) *discordgo.InteractionResponse {
	r, err := ctx.GetChannel(ctx.Guild.ID, channels.DISCORD)
	if err != nil {
		return ctx.Errorf("Room %s not active", ctx.Guild.ID)
	}
	r.Pause()
	return ctx.Reply(":pause_button: Pausing :thumbsup:")
}

func skipCMD(ctx CommandCtx) *discordgo.InteractionResponse {
	r, err := ctx.GetChannel(ctx.Guild.ID, channels.DISCORD)
	if err != nil {
		return ctx.Errorf("Room %s not active", ctx.Guild.ID)
	}
	r.Skip()
	return ctx.Reply(":fast_forward: ***Skipped*** :thumbsup:")
}
