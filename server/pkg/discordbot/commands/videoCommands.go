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
			Command:     "play",
			Description: "Plays a song with the given name or url",
			Usage:       "<link/query>",
			Aliases:     []string{"p", "add"},
			Function:    playCmd,
		},
		CMD{
			Command:     "pause",
			Description: "Pauses the current playing track",
			Aliases:     []string{"stop"},
			Function:    pauseCMD,
		},
		CMD{
			Command:     "watch",
			Description: "Get Url Watch2gether Room, where the video will be in sync with discord",
			Aliases:     []string{"link"},
			Function:    LinkCMD,
		},
		CMD{
			Command:     "skip",
			Description: "Skips the current song and plays the song you requested.",
			Aliases:     []string{"fs", "skipped", "next"},
			Function:    skipCMD,
		},
	)
}

func AddVideo(uri string, username string, meta *meta.Meta, r *room.Room) (*EmbededMessage, error) {

	u, err := url.ParseRequestURI(uri)
	if err != nil {
		return nil, fmt.Errorf("%s Is not a valid URL", uri)
	}
	videos := media.NewVideo(u.String(), username)
	video := videos[0]
	message := EmbedBuilder(fmt.Sprintf("Added %d tracks to the Queue", len(videos)))
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

	message.AddField(discordgo.MessageEmbedField{
		Name:   "Position in the Queue",
		Value:  fmt.Sprintf("%d", len(meta.Queue)+1),
		Inline: false,
	})
	queue := append(meta.Queue, videos...)
	r.HandleEvent(events.Event{
		Action:  events.EVNT_UPDATE_QUEUE,
		Watcher: user.DISCORD_BOT,
		Queue:   queue,
	})
	return message, nil
}

func playCmd(ctx CommandCtx) error {

	r, ok := ctx.GetHubRoom()
	meta, err := ctx.GetMeta()
	if !ok || err != nil {
		return fmt.Errorf("Room %s not active", ctx.Guild.ID)
	}

	if len(ctx.Args) > 0 {
		msg, err := AddVideo(ctx.Args[0], ctx.User.Username, meta, r)
		if err != nil {
			return err
		}
		ctx.ReplyEmbed(msg)
	}

	if !meta.Playing {
		evt := events.NewEvent(events.EVNT_PLAYING)
		evt.Watcher = user.DISCORD_BOT
		r.HandleEvent(evt)
		ctx.Reply(":play_pause: Resuming :thumbsup:")
	}
	return nil
}

func pauseCMD(ctx CommandCtx) error {
	r, ok := ctx.GetHubRoom()
	if !ok {
		return fmt.Errorf("Room %s not active", ctx.Guild.ID)
	}
	evt := events.NewEvent(events.EVNT_PAUSING)
	evt.Watcher = user.DISCORD_BOT
	r.HandleEvent(evt)
	return ctx.Reply(":pause_button: Pausing :thumbsup:")
}

func skipCMD(ctx CommandCtx) error {
	r, ok := ctx.GetHubRoom()
	if !ok {
		return fmt.Errorf("Room %s not active", ctx.Guild.ID)
	}
	r.HandleEvent(events.Event{
		Action:  events.EVNT_NEXT_VIDEO,
		Watcher: user.DISCORD_BOT,
	})
	return ctx.Reply(":fast_forward: ***Skipped*** :thumbsup:")
}
