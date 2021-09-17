package commands

import (
	"fmt"
	"strconv"
	"watch2gether/pkg/events"
	"watch2gether/pkg/media"
	"watch2gether/pkg/user"
)

func init() {
	Commands.Register(
		CMD{
			Command:     "move",
			Description: "Moves a certain song to a chosen position in the queue",
			Usage:       "!move <old positon> <new position>",
			Function:    moveCMD,
		},
		CMD{
			Command:     "remove",
			Description: "Removes a certain entry from the queue.",
			Aliases:     []string{"delete"},
			Usage:       "!remove <numbers>",
			Function:    removeCMD,
		},
		CMD{
			Command:     "skipTo",
			Description: "Skips to a certain position in the queue",
			Usage:       "!skipto <position>",
			Function:    skipToCMD,
		},
		CMD{
			Command:     "shuffle",
			Description: "Shuffles the entire queue",
			Function:    shuffleCMD,
		},
		CMD{
			Command:     "clear",
			Description: "clears the entire queue",
			Function:    clearCMD,
		},
	)
}

func shuffleCMD(ctx CommandCtx) error {
	r, ok := ctx.GetHubRoom()
	if !ok {
		return fmt.Errorf("room %s not active", ctx.Guild.ID)
	}
	r.HandleEvent(events.Event{
		Action:  events.EVT_SUFFLE_QUEUE,
		Watcher: user.DISCORD_BOT,
	})
	return ctx.Reply(":twisted_rightwards_arrows:  Queue Shuffled :thumbsup:")
}

func moveCMD(ctx CommandCtx) error {

	if len(ctx.Args) != 2 {
		return ctx.Reply(":cry: sorry not enough argunments in the command try `!move 1 2`")
	}
	pos1, err1 := strconv.Atoi(ctx.Args[0])
	pos2, err2 := strconv.Atoi(ctx.Args[1])

	//Convert user positions to slice locations
	pos1 = pos1 - 1
	pos2 = pos2 - 1

	if err1 != nil || err2 != nil {
		return ctx.Reply(":cry: sorry not enough argunments in the command try `!move 1 2`")
	}

	r, ok := ctx.GetHubRoom()
	meta, err := ctx.GetMeta()
	if !ok || err != nil {
		return fmt.Errorf("room %s not active", ctx.Guild.ID)
	}

	if pos1 < 0 || pos1 > len(meta.Queue)-1 || pos2 < 0 || pos2 > len(meta.Queue)-1 {
		return ctx.Reply(":cry: number not in range of the queue try again")
	}

	meta.ReorderQueue(pos1, pos2)
	r.HandleEvent(events.Event{
		Action:  events.EVNT_UPDATE_QUEUE,
		Watcher: user.DISCORD_BOT,
		Queue:   meta.Queue,
	})
	return ctx.Reply(":white_check_mark: Queue Updated :thumbsup:")
}

func removeCMD(ctx CommandCtx) error {

	if len(ctx.Args) > 1 {
		return ctx.Reply(":cry: sorry not enough argunments in the command try `!remove 1 2`")
	}
	r, ok := ctx.GetHubRoom()
	meta, err := ctx.GetMeta()
	if !ok || err != nil {
		return fmt.Errorf("room %s not active", ctx.Guild.ID)
	}

	for i, arg := range ctx.Args {
		pos, err := strconv.Atoi(arg)
		if err != nil {
			break
		}
		pos = pos - 1
		if i != 0 {
			pos = pos - i
		}
		if pos < 0 || pos > len(meta.Queue)-1 {
			continue
		}
		meta.RemoveFromQueue(pos)
	}
	r.HandleEvent(events.Event{
		Action:  events.EVNT_UPDATE_QUEUE,
		Watcher: user.DISCORD_BOT,
		Queue:   meta.Queue,
	})
	return ctx.Reply(":white_check_mark: Queue Updated :thumbsup:")
}

func skipToCMD(ctx CommandCtx) error {

	if len(ctx.Args) == 1 {
		ctx.Reply(":cry: sorry not enough argunments in the command try `!skipto 10`")
	}
	pos, err := strconv.Atoi(ctx.Args[0])
	if err != nil {
		ctx.Reply(":cry: sorry not enough argunments in the command try `!sktoTo 10`")
	}

	r, ok := ctx.GetHubRoom()
	meta, err := ctx.GetMeta()
	if !ok || err != nil {
		return fmt.Errorf("room %s not active", ctx.Guild.ID)
	}

	if pos < 0 || pos > len(meta.Queue)-1 {
		return ctx.Reply(":cry: number not in range of the queue try again")
	}
	for i := 0; i < pos; i++ {
		meta.RemoveFromQueue(0)
	}

	r.HandleEvent(events.Event{
		Action:  events.EVNT_UPDATE_QUEUE,
		Watcher: user.DISCORD_BOT,
		Queue:   meta.Queue,
	})
	r.HandleEvent(events.Event{
		Action:  events.EVNT_NEXT_VIDEO,
		Watcher: user.DISCORD_BOT,
	})

	return ctx.Reply(":white_check_mark: Skiped :thumbsup:")
}

func clearCMD(ctx CommandCtx) error {

	r, ok := ctx.GetHubRoom()
	if !ok {
		return fmt.Errorf("room %s not active", ctx.Guild.ID)
	}

	r.HandleEvent(events.Event{
		Action:  events.EVNT_UPDATE_QUEUE,
		Watcher: user.DISCORD_BOT,
		Queue:   []media.Video{},
	})

	return ctx.Reply(":white_check_mark: Cleared Queue :thumbsup:")
}
