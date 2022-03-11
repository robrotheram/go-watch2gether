package commands

import (
	"fmt"
	"net/url"
	"watch2gether/pkg/events"
	"watch2gether/pkg/media"
	"watch2gether/pkg/user"

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

func AddVideo(ctx CommandCtx, uri string) {

	r, _ := ctx.GetHubRoom()
	meta, _ := ctx.GetMeta()

	videos, err := media.NewVideo(uri, ctx.User.User.Username)
	if len(videos) == 0 || err != nil {
		return
	}
	video := videos[0]
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
		Value:  fmt.Sprintf("%d", len(meta.Queue)+1),
		Inline: false,
	})
	queue := append(meta.Queue, videos...)
	r.HandleEvent(events.Event{
		Action:  events.EVENT_UPDATE_QUEUE,
		Watcher: user.DISCORD_BOT,
		Queue:   queue,
	})
	ctx.ReplyEmbed(message)

	meta, _ = ctx.GetMeta()
	if meta.CurrentVideo.Url == "" {
		r.HandleEvent(events.Event{
			Action:  events.EVENT_NEXT_VIDEO,
			Watcher: user.DISCORD_BOT,
		})
		ctx.Reply(":play_pause: Now Playing :thumbsup:")
	}

}

func addCmd(ctx CommandCtx) *discordgo.InteractionResponse {

	if _, ok := ctx.GetHubRoom(); !ok {
		return ctx.Errorf("room %s not active please get bot to join room", ctx.Guild.ID)
	}

	if len(ctx.Args) <= 0 {
		return ctx.Errorf("no media url found")
	}

	_, err := url.ParseRequestURI(ctx.Args[0])
	if err != nil {
		return ctx.Errorf("%s Is not a valid URL", ctx.Args[0])
	}

	if media.MediaFactory.GetFactory(ctx.Args[0]) == nil {
		return ctx.Errorf("%s is a unsupported media type", ctx.Args[0])
	}

	go AddVideo(ctx, ctx.Args[0])

	return nil
}

func playCmd(ctx CommandCtx) *discordgo.InteractionResponse {
	r, ok := ctx.GetHubRoom()
	if !ok {
		return ctx.Errorf("Room %s not active", ctx.Guild.ID)
	}
	r.HandleEvent(events.Event{
		Action:  events.EVENT_PLAYING,
		Watcher: user.DISCORD_BOT,
	})
	return ctx.Reply(":play_pause: Now Playing :thumbsup:")
}

func pauseCMD(ctx CommandCtx) *discordgo.InteractionResponse {
	r, ok := ctx.GetHubRoom()
	if !ok {
		return ctx.Errorf("Room %s not active", ctx.Guild.ID)
	}
	evt := events.NewEvent(events.EVENT_PAUSING)
	evt.Watcher = user.DISCORD_BOT
	r.HandleEvent(evt)
	return ctx.Reply(":pause_button: Pausing :thumbsup:")
}

func skipCMD(ctx CommandCtx) *discordgo.InteractionResponse {
	r, ok := ctx.GetHubRoom()
	if !ok {
		return ctx.Errorf("Room %s not active", ctx.Guild.ID)
	}
	r.HandleEvent(events.Event{
		Action:  events.EVENT_NEXT_VIDEO,
		Watcher: user.DISCORD_BOT,
	})
	return ctx.Reply(":fast_forward: ***Skipped*** :thumbsup:")
}
