package commands

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

func init() {
	Commands.Register(
		CMD{
			Command:     "help",
			Description: "",
			Function: func(ctx CommandCtx) error {
				embed := EmbedBuilder("Watch2gether Help")
				msg := ""

				sortkey := Commands.SortKeys()

				for _, v := range sortkey {
					cmd := Commands.Cmds[v]
					msg = msg + cmd.Format()
				}
				fmt.Println(len(msg))
				embed.AddField(discordgo.MessageEmbedField{
					Name:   "Avalible Commands",
					Value:  msg,
					Inline: true,
				})
				return ctx.ReplyEmbed(embed)
			},
		})
}
