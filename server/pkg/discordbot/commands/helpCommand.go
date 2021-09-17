package commands

import "github.com/bwmarrin/discordgo"

func init() {
	Register(
		CMD{
			Command:     "help",
			Description: "",
			Function: func(ctx CommandCtx) error {
				embed := EmbedBuilder("Watch2gether Help")
				msg := ""

				sortkey := SortKeys(Cmds)

				for _, v := range sortkey {
					cmd := Cmds[v]
					msg = msg + cmd.Format()
				}
				embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
					Name:   "Avalible Commands",
					Value:  msg,
					Inline: true,
				})
				return ctx.ReplyEmbed(embed)
			},
		})
}
