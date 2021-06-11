package command

import (
	"fmt"
	"net/url"
	"watch2gether/pkg/events"
	"watch2gether/pkg/media"
	"watch2gether/pkg/user"
	"watch2gether/pkg/utils"

	"github.com/segmentio/ksuid"
)

func init() {
	Commands["add"] = &AddCmd{BaseCommand{"Add Video to Queue"}}
	Commands["status"] = &StatusCmd{BaseCommand{"Current Status of what is playing"}}
	Commands["skip"] = &SkipCmd{BaseCommand{"Skip to next video in the Queue"}}
	Commands["queue"] = &listCmd{BaseCommand{"List videos in the Queue"}}
	Commands["shuffle"] = &SuffleCmd{BaseCommand{"Shuffle the Queue"}}
}

type AddCmd struct{ BaseCommand }

func (cmd *AddCmd) Execute(ctx CommandCtx) error {
	u, err := url.ParseRequestURI(ctx.Args[0])
	if err != nil {
		return fmt.Errorf("%s Is not a valid URL", ctx.Args[0])
	}
	r, ok := ctx.GetHubRoom()
	meta, err := ctx.GetMeta()
	if !ok && err != nil {
		return fmt.Errorf("Room %s not active", ctx.Guild.ID)
	}
	document, err := utils.Scrape(u.String(), 1)
	if err != nil {
		return fmt.Errorf("Video Error not found")
	}
	video := media.Video{ID: ksuid.New().String(), Title: document.Preview.Title, Url: u.String(), User: ctx.User.Username}
	queue := append(meta.Queue, video)
	r.HandleEvent(events.Event{
		Action:  events.EVNT_UPDATE_QUEUE,
		Watcher: user.DISCORD_BOT,
		Queue:   queue,
	})
	return ctx.Reply(fmt.Sprintf("Video %s added to room Queue", video.Title))
}

type StatusCmd struct{ BaseCommand }

func (cmd *StatusCmd) Execute(ctx CommandCtx) error {
	meta, err := ctx.GetMeta()
	if err != nil {
		return fmt.Errorf("Room %s not active", ctx.Guild.ID)
	}
	vidoe := meta.CurrentVideo
	if vidoe.Url == "" {
		return ctx.Reply(fmt.Sprintf("No Video Playing"))
	}
	return ctx.Reply(fmt.Sprintf("Currently Playing: %s \n 📼 %s ", vidoe.Title, vidoe.Url))
}

type SkipCmd struct{ BaseCommand }

func (cmd *SkipCmd) Execute(ctx CommandCtx) error {
	r, ok := ctx.GetHubRoom()
	if !ok {
		return fmt.Errorf("Room %s not active", ctx.Guild.ID)
	}
	r.HandleEvent(events.Event{
		Action:  events.EVNT_NEXT_VIDEO,
		Watcher: user.DISCORD_BOT,
	})
	return ctx.Reply("Video Skipped")
}

type SuffleCmd struct{ BaseCommand }

func (cmd *SuffleCmd) Execute(ctx CommandCtx) error {
	r, ok := ctx.GetHubRoom()
	if !ok {
		return fmt.Errorf("Room %s not active", ctx.Guild.ID)
	}
	r.HandleEvent(events.Event{
		Action:  events.EVT_SUFFLE_QUEUE,
		Watcher: user.DISCORD_BOT,
	})
	return ctx.Reply("Queue Shuffled")
}

type listCmd struct{ BaseCommand }

func (cmd *listCmd) Execute(ctx CommandCtx) error {
	meta, err := ctx.GetMeta()
	if err != nil {
		return fmt.Errorf("Room %s not active", ctx.Guild.ID)
	}
	msg := ""
	vidoe := meta.CurrentVideo
	if vidoe.Url != "" {
		msg += fmt.Sprintf("Currently Playing: \n %s \n", vidoe.Title)
		msg += "--------------------------------------------------\n "
	}
	msg += "Queue: \n"
	for i, v := range meta.Queue {
		msg = msg + fmt.Sprintf(" -%d: %s \n", i+1, v.Title)
	}
	return ctx.Reply(msg)
}
