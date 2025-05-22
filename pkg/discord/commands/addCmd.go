package commands

import (
	"fmt"
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

	go func(ctx CommandCtx) {
		isTop := getPostionOption(ctx) == "TOP"
		url := ctx.Args[0]
		// This runs in the background
		err = ctx.Controller.Add(url, isTop, ctx.Member.User.Username)
		var content string
		if err != nil {
			content = fmt.Sprintf("Failed to add video: %v", err)
		} else {
			content = "Video successfully added to the channel!"
		}
		// Send follow-up message
		ctx.Session.FollowupMessageCreate(ctx.Interaction, false, &discordgo.WebhookParams{
			Content: content,
		})
	}(ctx)

	return ctx.Replyf("⏳ Processing your media link: `%s`\nIt will soon be added to the queue give it a mo", ctx.Args[0])
}
