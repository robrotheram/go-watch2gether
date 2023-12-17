package commands

import (
	"strconv"

	"github.com/bwmarrin/discordgo"
)

var integerOptionMinValue = 0.0

func init() {
	register(
		Command{
			Name: "move",
			ApplicationCommand: []discordgo.ApplicationCommand{
				{

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
			},
			Function: moveCMD,
		},
		Command{
			Name: "remove",
			ApplicationCommand: []discordgo.ApplicationCommand{
				{
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
			},
			Usage:    "!remove <numbers>",
			Function: removeCMD,
		},
		Command{
			Name: "shuffle",
			ApplicationCommand: []discordgo.ApplicationCommand{
				{
					Description: "Shuffles the entire queue",
				},
			},
			Function: shuffleCMD,
		},
		Command{
			Name: "clear",
			ApplicationCommand: []discordgo.ApplicationCommand{
				{
					Description: "clears the entire queue",
				},
			},
			Function: clearCMD,
		},
	)
}

func shuffleCMD(ctx CommandCtx) *discordgo.InteractionResponse {
	state := ctx.Controller.State()
	state.Shuffle()
	ctx.Controller.Update(state)
	return ctx.Reply(":twisted_rightwards_arrows:  Queue Shuffled :thumbsup:")
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
	state := ctx.Controller.State()
	if pos1 < 0 || pos1 > len(state.Queue)-1 || pos2 < 0 || pos2 > len(state.Queue)-1 {
		return ctx.Reply(":cry: number not in range of the queue try again")
	}
	state.Reorder(pos1, pos2)

	return ctx.Reply(":white_check_mark: Queue Updated :thumbsup:")
}

func removeCMD(ctx CommandCtx) *discordgo.InteractionResponse {
	state := ctx.Controller.State()
	if len(ctx.Args) > 1 {
		return ctx.Reply(":cry: sorry not enough argunments in the command try `!remove 1 2`")
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
		if pos < 0 || pos > len(state.Queue)-1 {
			continue
		}
		state.Remove(pos)
	}
	ctx.Controller.Update(state)
	return ctx.Reply(":white_check_mark: Queue Updated :thumbsup:")
}

func clearCMD(ctx CommandCtx) *discordgo.InteractionResponse {
	state := ctx.Controller.State()
	state.Clear()
	ctx.Controller.Update(state)
	return ctx.Reply(":white_check_mark: Cleared Queue :thumbsup:")
}
