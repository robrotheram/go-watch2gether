package commands

import (
	"fmt"
	"net/url"
	"watch2gether/pkg/events"
	"watch2gether/pkg/media"
	"watch2gether/pkg/room"
	meta "watch2gether/pkg/roomMeta"
	"watch2gether/pkg/user"

	"github.com/bwmarrin/discordgo"
)

func init() {
	Commands.Register(
		CMD{
			ApplicationCommand: discordgo.ApplicationCommand{
				Name:        "play",
				Description: "Plays a song with the given name or url",
				Options: []*discordgo.ApplicationCommandOption{
					{
						Type:        discordgo.ApplicationCommandOptionString,
						Name:        "video",
						Description: "Video URL e.g (https://www.youtube.com/watch?v=noneMROp_E8)",
						Required:    false,
					},
				},
			},

			Function: playCmd,
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

func AddVideo(uri string, username string, meta *meta.Meta, r *room.Room) (*EmbededMessage, error) {

	u, err := url.ParseRequestURI(uri)
	if err != nil {
		return nil, fmt.Errorf("%s Is not a valid URL", uri)
	}
	videos, err := media.NewVideo(u.String(), username)
	if len(videos) == 0 || err != nil {
		return nil, fmt.Errorf("unable to understand the video does it exist?")
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
	fmt.Println(message)
	return message, nil
}

func playCmd(ctx CommandCtx) *discordgo.InteractionResponse {

	r, ok := ctx.GetHubRoom()
	meta, err := ctx.GetMeta()
	if !ok || err != nil {
		return ctx.Errorf("Room %s not active", ctx.Guild.ID)
	}

	if len(ctx.Args) > 0 {
		msg, err := AddVideo(ctx.Args[0], ctx.User.User.Username, meta, r)
		if err != nil {
			return ctx.Errorf("Unable to add video: %v", err)
		}
		ctx.ReplyEmbed(msg)
	} else {
		r.HandleEvent(events.Event{
			Action:  events.EVENT_PLAYING,
			Watcher: user.DISCORD_BOT,
		})
		return ctx.Reply(":play_pause: Resuming :thumbsup:")
	}
	meta, _ = ctx.GetMeta()
	if meta.CurrentVideo.Url == "" {
		r.HandleEvent(events.Event{
			Action:  events.EVENT_NEXT_VIDEO,
			Watcher: user.DISCORD_BOT,
		})
		ctx.Reply(":play_pause: Now Playing :thumbsup:")
	}
	return nil
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
