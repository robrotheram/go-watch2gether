package commands

import (
	"fmt"
	"w2g/pkg/discord/players"

	"github.com/bwmarrin/discordgo"
)

func init() {
	register(
		Command{
			ApplicationCommand: discordgo.ApplicationCommand{
				Name: "join",
				Type: discordgo.UserApplicationCommand,
			},
			Function: joinCmd,
		},
		Command{
			ApplicationCommand: discordgo.ApplicationCommand{
				Name: "join",
				Type: discordgo.ChatApplicationCommand,
			},
			Function: joinCmd,
		}, Command{
			ApplicationCommand: discordgo.ApplicationCommand{
				Name: "leave",
				Type: discordgo.UserApplicationCommand,
			},
			Function: leave,
		},
		Command{
			ApplicationCommand: discordgo.ApplicationCommand{
				Name: "leave",
				Type: discordgo.ChatApplicationCommand,
			},
			Function: leave,
		})

}

func GetUserVoiceChannel(session *discordgo.Session, user string) (string, error) {
	for _, g := range session.State.Guilds {
		for _, v := range g.VoiceStates {
			if v.UserID == user {
				return v.ChannelID, nil
			}
		}
	}
	return "", fmt.Errorf("channel not found")
}

func join(ctx CommandCtx) error {
	vc, err := GetUserVoiceChannel(ctx.Session, ctx.User.User.ID)
	if err != nil {
		return fmt.Errorf("you are not connected to voice channel")
	}
	voice, err := ctx.Session.ChannelVoiceJoin(ctx.Guild.ID, vc, false, true)
	if err != nil || voice == nil {
		return fmt.Errorf("you are not connected to voice channel")
	}
	ctx.Controller.Join(players.NewDiscordPlayer(voice))
	return nil
}

func joinCmd(ctx CommandCtx) *discordgo.InteractionResponse {
	if join(ctx) != nil {
		return ctx.Reply("User not connected to voice channel")
	}
	return ctx.Reply("🔉 Connected to voice channel 🔉")
}

func leave(ctx CommandCtx) *discordgo.InteractionResponse {
	ctx.Controller.Leave()
	return ctx.Reply("👋 cheerio have a good day 🎩")
}
