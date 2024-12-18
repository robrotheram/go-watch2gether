package commands

import (
	"log"
	"net/url"

	"github.com/bwmarrin/discordgo"
)

func init() {
	register(Command{
		Name: "add",
		ApplicationCommand: []discordgo.ApplicationCommand{
			{
				Description: "Add new track to the queue",
				Options: []*discordgo.ApplicationCommandOption{
					{
						Type:        discordgo.ApplicationCommandOptionString,
						Name:        "media",
						Description: "Video/Audio URL e.g (https://www.youtube.com/watch?v=noneMROp_E8)",
						Required:    true,
					},
					{
						Type:        discordgo.ApplicationCommandOptionString,
						Name:        "position",
						Description: "Where to add new media top of bottom of the queue (default: bottom)",
						Required:    false,
						Choices: []*discordgo.ApplicationCommandOptionChoice{
							{
								Name:  "top of queue",
								Value: "TOP",
							},
							{
								Name:  "bottom of queue",
								Value: "BOTTOM",
							},
						},
					},
				},
				Type: discordgo.ChatApplicationCommand,
			},
		},
		Function: addCmd,
	})
}

func addCmd(ctx CommandCtx) *discordgo.InteractionResponse {
	_, err := url.ParseRequestURI(ctx.Args[0])
	if err != nil {
		return ctx.Errorf("%s Is not a valid URL", ctx.Args[0])
	}
	isTop := getPostionOption(ctx) == "TOP"
	go func() {
		err = ctx.Controller.Add(ctx.Args[0], isTop, ctx.Member.User.Username)
		if err != nil {
			//ctx.Errorf("error: %v", err)
			log.Printf("Unable to add vidoe: %v", err)
			return
		}

		// if !ctx.Controller.ContainsPlayer(ctx.Guild.ID) {
		// 	join(ctx)
		// }

		// if ctx.Controller.State().State != controllers.PLAY {
		// 	ctx.Controller.Start(ctx.Member.User.Username)
		// }
	}()

	return ctx.Replyf("⏳ Processing your media link: `%s`\nIt will soon be added to the queue give it a mo", ctx.Args[0])
}
