package discordbot

import (
	"strconv"

	"github.com/bwmarrin/discordgo"
)

var integerOptionMinValue = 0.0

func init() {
	Commands.Register(
		CMD{
			ApplicationCommand: discordgo.ApplicationCommand{
				Name:        "move",
				Description: "Moves a certain song to a chosen position in the queue",
				Options: []*discordgo.ApplicationCommandOption{
					{
						Type:        discordgo.ApplicationCommandOptionInteger,
						Name:        "old-position",
						Description: "Current position in the queue",
						Required:    true,
					},
					{
						Type:        discordgo.ApplicationCommandOptionInteger,
						Name:        "new-position",
						Description: "New position in the queue",
						Required:    true,
					},
				},
			},
			Function: moveCMD,
		},
		CMD{
			ApplicationCommand: discordgo.ApplicationCommand{
				Name:        "remove",
				Description: "Removes a certain entry from the queue.",
				Options: []*discordgo.ApplicationCommandOption{
					{
						Type:        discordgo.ApplicationCommandOptionInteger,
						Name:        "remove-position",
						Description: "Remove current position in the queue",
						MinValue:    &integerOptionMinValue,
						Required:    true,
					},
				},
			},
			Usage:    "!remove <numbers>",
			Function: removeCMD,
		},
		CMD{
			ApplicationCommand: discordgo.ApplicationCommand{
				Name:        "skip",
				Description: "Skips to a certain position in the queue",
				Options: []*discordgo.ApplicationCommandOption{
					{
						Type:        discordgo.ApplicationCommandOptionInteger,
						Name:        "skip-position",
						Description: "Skip to this poisition in the queue",
						MinValue:    &integerOptionMinValue,
						Required:    true,
					},
				},
			},
			Usage:    "!skip <position>",
			Function: skipToCMD,
		},
		CMD{
			ApplicationCommand: discordgo.ApplicationCommand{
				Name:        "shuffle",
				Description: "Shuffles the entire queue",
			},
			Function: shuffleCMD,
		},
		CMD{
			ApplicationCommand: discordgo.ApplicationCommand{
				Name:        "loop",
				Description: "enable/disable loop the current video",
			},
			Function: loopCMD,
		},
		CMD{
			ApplicationCommand: discordgo.ApplicationCommand{
				Name:        "clear",
				Description: "clears the entire queue",
			},
			Function: clearCMD,
		},
	)
}

func shuffleCMD(ctx CommandCtx) *discordgo.InteractionResponse {
	r, err := ctx.GetChannel(ctx.Guild.ID)
	if err != nil {
		return ctx.Errorf("room %s not active", ctx.Guild.ID)
	}
	r.Shuffle()
	return ctx.Reply(":twisted_rightwards_arrows:  Queue Shuffled :thumbsup:")
}

func loopCMD(ctx CommandCtx) *discordgo.InteractionResponse {
	r, err := ctx.GetChannel(ctx.Guild.ID)
	if err != nil {
		return ctx.Errorf("room %s not active", ctx.Guild.ID)
	}
	state, _ := r.GetState()
	r.SetLoop(state.Loop)
	if state.Loop {
		return ctx.Reply(":arrows_counterclockwise:  Looping enabled")
	}
	return ctx.Reply(":arrows_counterclockwise:  Looping disabled")
}

func moveCMD(ctx CommandCtx) *discordgo.InteractionResponse {
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

	r, err := ctx.GetChannel(ctx.Guild.ID)
	if err != nil {
		return ctx.Errorf("room %s not active", ctx.Guild.ID)
	}
	queue, _ := r.GetQueue()
	if pos1 < 0 || pos1 > len(queue)-1 || pos2 < 0 || pos2 > len(queue)-1 {
		return ctx.Reply(":cry: number not in range of the queue try again")
	}

	r.Move(pos1, pos2)

	return ctx.Reply(":white_check_mark: Queue Updated :thumbsup:")
}

func removeCMD(ctx CommandCtx) *discordgo.InteractionResponse {

	if len(ctx.Args) > 1 {
		return ctx.Reply(":cry: sorry not enough argunments in the command try `!remove 1 2`")
	}
	r, err := ctx.GetChannel(ctx.Guild.ID)
	if err != nil {
		return ctx.Errorf("room %s not active", ctx.Guild.ID)
	}
	queue, _ := r.GetQueue()
	for i, arg := range ctx.Args {
		pos, err := strconv.Atoi(arg)
		if err != nil {
			break
		}
		pos = pos - 1
		if i != 0 {
			pos = pos - i
		}
		if pos < 0 || pos > len(queue)-1 {
			continue
		}
		r.Remove(pos)
	}
	return ctx.Reply(":white_check_mark: Queue Updated :thumbsup:")
}

func skipToCMD(ctx CommandCtx) *discordgo.InteractionResponse {

	if len(ctx.Args) == 1 {
		ctx.Reply(":cry: sorry not enough argunments in the command try `!skipto 10`")
	}
	pos, err := strconv.Atoi(ctx.Args[0])
	if err != nil {
		ctx.Reply(":cry: sorry not enough argunments in the command try `!sktoTo 10`")
	}

	r, err := ctx.GetChannel(ctx.Guild.ID)
	if err != nil {
		return ctx.Errorf("Room %s not active", ctx.Guild.ID)
	}
	queue, _ := r.GetQueue()
	if pos < 0 || pos > len(queue)-1 {
		return ctx.Reply(":cry: number not in range of the queue try again")
	}
	for i := 0; i < pos; i++ {
		r.Remove(pos)
	}
	r.Skip()
	return ctx.Reply(":white_check_mark: Skiped :thumbsup:")
}

func clearCMD(ctx CommandCtx) *discordgo.InteractionResponse {
	r, err := ctx.GetChannel(ctx.Guild.ID)
	if err != nil {
		return ctx.Errorf("Room %s not active", ctx.Guild.ID)
	}
	r.Clear()
	return ctx.Reply(":white_check_mark: Cleared Queue :thumbsup:")
}