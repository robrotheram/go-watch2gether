package commands

import (
	"fmt"
	"strconv"
	"watch2gether/pkg/events"
	"watch2gether/pkg/user"
)

func init() {
	Register(
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
		},
		CMD{
			Command:     "shuffle",
			Description: "Shuffles the entire queue",
			Function:    shuffleCMD,
		},
	)
}

func shuffleCMD(ctx CommandCtx) error {
	r, ok := ctx.GetHubRoom()
	if !ok {
		return fmt.Errorf("Room %s not active", ctx.Guild.ID)
	}
	r.HandleEvent(events.Event{
		Action:  events.EVT_SUFFLE_QUEUE,
		Watcher: user.DISCORD_BOT,
	})
	return ctx.Reply(":twisted_rightwards_arrows:  Queue Shuffled :thumbsup:")
}

func moveCMD(ctx CommandCtx) error {

	if len(ctx.Args) != 2 {
		fmt.Errorf(":cry: sorry not enough argunments in the command try `!move 1 2`")
	}
	pos1, err1 := strconv.Atoi(ctx.Args[0])
	pos2, err2 := strconv.Atoi(ctx.Args[1])

	if err1 != nil || err2 != nil {
		fmt.Errorf(":cry: sorry not enough argunments in the command try `!move 1 2`")
	}

	r, ok := ctx.GetHubRoom()
	meta, err := ctx.GetMeta()
	if !ok || err != nil {
		return fmt.Errorf("Room %s not active", ctx.Guild.ID)
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
		fmt.Errorf(":cry: sorry not enough argunments in the command try `!remove 1 2`")
	}
	r, ok := ctx.GetHubRoom()
	meta, err := ctx.GetMeta()
	if !ok || err != nil {
		return fmt.Errorf("Room %s not active", ctx.Guild.ID)
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
		fmt.Errorf(":cry: sorry not enough argunments in the command try `!skipto 10`")
	}
	pos, err := strconv.Atoi(ctx.Args[0])
	if err != nil {
		fmt.Errorf(":cry: sorry not enough argunments in the command try `!sktoTo 10`")
	}

	r, ok := ctx.GetHubRoom()
	meta, err := ctx.GetMeta()
	if !ok || err != nil {
		return fmt.Errorf("Room %s not active", ctx.Guild.ID)
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
