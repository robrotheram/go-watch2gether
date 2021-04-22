package command

import (
	"fmt"
	"watch2gether/pkg/events"
	"watch2gether/pkg/user"
)

func init() {
	Commands["play"] = &PlayCmd{BaseCommand{"Play Video"}}
	Commands["pause"] = &PauseCmd{BaseCommand{"Pause Video"}}
	Commands["stop"] = &PauseCmd{BaseCommand{"Stop Video"}}
	Commands["watch"] = &WatchCmd{BaseCommand{"Get Link to room"}}
}

type PlayCmd struct{ BaseCommand }
type PauseCmd struct{ BaseCommand }
type WatchCmd struct{ BaseCommand }

func (cmd *PlayCmd) Execute(ctx CommandCtx) error {
	r, ok := ctx.GetHubRoom()
	if !ok {
		return fmt.Errorf("Room %s not active", ctx.Guild.ID)
	}
	evt := events.NewEvent(events.EVNT_PLAYING)
	evt.Watcher = user.DISCORD_BOT
	r.HandleEvent(evt)
	return nil
}

func (cmd *PauseCmd) Execute(ctx CommandCtx) error {
	r, ok := ctx.GetHubRoom()
	if !ok {
		return fmt.Errorf("Room %s not active", ctx.Guild.ID)
	}
	evt := events.NewEvent(events.EVNT_PAUSING)
	evt.Watcher = user.DISCORD_BOT
	r.HandleEvent(evt)
	return nil
}

func (cmd *WatchCmd) Execute(ctx CommandCtx) error {
	ctx.Reply(fmt.Sprintf("📼 %s/app/room/%s 📼 \n Stay Awhile and Listen", ctx.BaseURL, ctx.Guild.ID))
	return nil
}
