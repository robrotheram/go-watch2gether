package commands

import (
	"fmt"
	"w2g/pkg/discord/players"

	"github.com/bwmarrin/discordgo"
)

func init() {
	register(
		Command{
			Name: "join",
			ApplicationCommand: []discordgo.ApplicationCommand{
				{
					Type: discordgo.UserApplicationCommand,
				},
				{
					Description: "Join the voice channle you are in",
					Type:        discordgo.ChatApplicationCommand,
				},
			},
			Function: joinCmd,
		},
		Command{
			Name: "leave",
			ApplicationCommand: []discordgo.ApplicationCommand{
				{
					Type: discordgo.UserApplicationCommand,
				},
				{
					Description: "Leave the voice channle you are in",
					Type:        discordgo.ChatApplicationCommand,
				},
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
	vc, err := GetUserVoiceChannel(ctx.Session, ctx.Member.User.ID)
	if err != nil {
		return fmt.Errorf("you are not connected to voice channel")
	}
	voice, err := ctx.Session.ChannelVoiceJoin(ctx.Guild.ID, vc, false, true)
	if err != nil || voice == nil {
		return fmt.Errorf("you are not connected to voice channel")
	}
	ctx.Controller.Stop(ctx.Member.User.Username)
	ctx.Controller.Join(players.NewDiscordPlayer(voice), ctx.Member.User.Username)
	return nil
}

func joinCmd(ctx CommandCtx) *discordgo.InteractionResponse {
	if join(ctx) != nil {
		return ctx.Reply("User not connected to voice channel")
	}
	return ctx.Reply("🔉 Connected to voice channel 🔉")
}

func leave(ctx CommandCtx) *discordgo.InteractionResponse {
	ctx.Controller.Leave(players.DISCORD, ctx.Member.User.Username)
	return ctx.Reply("👋 cheerio have a good day 🎩")
}
