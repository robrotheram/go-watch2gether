package command

import (
	"fmt"
	"net/url"
	"watch2gether/pkg/media"
	"watch2gether/pkg/user"
	"watch2gether/pkg/utils"

	"github.com/segmentio/ksuid"
)

func init() {
	Commands["add"] = &AddCmd{BaseCommand{"Add Video to Queue"}}
	Commands["status"] = &StatusCmd{BaseCommand{"Current Status of what is playing"}}
	Commands["skip"] = &SkipCmd{BaseCommand{"Skip to next video in the Queue"}}
}

type AddCmd struct{ BaseCommand }

func (cmd *AddCmd) Execute(ctx CommandCtx) error {
	u, err := url.ParseRequestURI(ctx.Args[0])
	if err != nil {
		return fmt.Errorf("%s Is not a valid URL", ctx.Args[0])
	}
	r, ok := ctx.GetHubRoom()
	if !ok {
		return fmt.Errorf("Room %s not active", ctx.Guild.ID)
	}
	document, err := utils.Scrape(u.String(), 1)
	if err != nil {
		return fmt.Errorf("Video Error not found")
	}
	video := media.Video{ID: ksuid.New().String(), Title: document.Preview.Title, Url: u.String(), User: ctx.User.Username}
	r.AddVideo(video, user.DISCORD_BOT)
	return ctx.Reply(fmt.Sprintf("Video %s added to room Queue", video.Title))
}

type StatusCmd struct{ BaseCommand }

func (cmd *StatusCmd) Execute(ctx CommandCtx) error {
	r, ok := ctx.GetHubRoom()
	if !ok {
		return fmt.Errorf("Room %s not active", ctx.Guild.ID)
	}
	vidoe := r.GetVideo()
	if vidoe.ID == "" {
		return ctx.Reply(fmt.Sprintf("No Video Playing"))
	}
	return ctx.Reply(fmt.Sprintf("Currently Playing: %s", vidoe.Title))
}

type SkipCmd struct{ BaseCommand }

func (cmd *SkipCmd) Execute(ctx CommandCtx) error {
	r, ok := ctx.GetHubRoom()
	if !ok {
		return fmt.Errorf("Room %s not active", ctx.Guild.ID)
	}

	r.ChangeVideo(user.DISCORD_BOT)
	return ctx.Reply("Video Skipped")
}
