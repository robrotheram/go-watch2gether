package discordbot

import (
	"fmt"
	"watch2gether/pkg/players"

	"github.com/bwmarrin/discordgo"
)

func init() {
	Commands.Register(
		CMD{
			ApplicationCommand: discordgo.ApplicationCommand{
				Name:        "join",
				Description: "Summons the bot to the voice channel you are in",
			},
			Function: JoinCmd,
		},

		CMD{
			ApplicationCommand: discordgo.ApplicationCommand{
				Name:        "leave",
				Description: "Disconnects the bot from the voice channel it is in.",
			},
			Function: LeaveCmd,
		},
	)
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

func JoinCmd(ctx CommandCtx) *discordgo.InteractionResponse {
	vc, err := GetUserVoiceChannel(ctx.Session, ctx.User.User.ID)
	if err != nil {
		ctx.Reply("User not connected to voice channel")
	}
	voice, err := ctx.Session.ChannelVoiceJoin(ctx.Guild.ID, vc, false, true)
	if err != nil {
		ctx.Reply("User not connected to voice channel")
	}
	conroller := players.NewDiscordPlayer(ctx.Guild.ID, voice)
	ctx.RegisterNewChannel(ctx.Guild.ID, conroller)
	return ctx.Reply("Bot added to the room")
}

func LeaveCmd(ctx CommandCtx) *discordgo.InteractionResponse {
	err := ctx.LeaveChannel(ctx.Guild.ID)
	if err != nil {
		ctx.Reply("Error Bot not connected")
	}
	return nil
}
