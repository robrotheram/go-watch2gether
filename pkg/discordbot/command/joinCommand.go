package command

import (
	"fmt"
	"watch2gether/pkg/audioBot"
	"watch2gether/pkg/user"

	"github.com/bwmarrin/discordgo"
)

func init() {
	Commands["join"] = &JoinCmd{BaseCommand{"Join Bot to a channel"}}
	Commands["leave"] = &LeaveCmd{BaseCommand{"Disconnect Bot from channel"}}
}

type JoinCmd struct{ BaseCommand }
type LeaveCmd struct{ BaseCommand }

func GetUserVoiceChannel(session *discordgo.Session, user string) (string, error) {
	for _, g := range session.State.Guilds {
		for _, v := range g.VoiceStates {
			if v.UserID == user {
				return v.ChannelID, nil
			}
		}
	}
	return "", fmt.Errorf("Channel Not found")
}

func (cmd *JoinCmd) Execute(ctx CommandCtx) error {
	r, ok := ctx.GetHubRoom()
	if !ok {
		return fmt.Errorf("Room %s not active", ctx.Guild.ID)
	}
	vc, err := GetUserVoiceChannel(ctx.Session, ctx.User.ID)
	if err != nil {
		ctx.Reply("User not connected to voice channel")
	}
	voice, err := ctx.Session.ChannelVoiceJoin(ctx.Guild.ID, vc, false, true)
	if err != nil {
		ctx.Reply("User not connected to voice channel")
	}
	bot := audioBot.NewAudioBot("", ctx.Channel.ID, voice, ctx.Session)
	err = r.RegisterBot(bot)
	if err != nil {
		ctx.Reply(fmt.Sprintf("Bot error %v", err))
	}
	ctx.Reply(fmt.Sprintf("Bot added to voice channel"))
	//_, err = ctx.Session.ChannelVoiceJoin(ctx.Guild.ID, vc, false, true)
	return nil
}

func (cmd *LeaveCmd) Execute(ctx CommandCtx) error {
	r, ok := ctx.GetHubRoom()
	if !ok {
		return fmt.Errorf("Room %s not active", ctx.Guild.ID)
	}
	r.Leave(user.DISCORD_BOT.ID)
	if r.Bot != nil {
		return r.Bot.Disconnect()
	}
	return ctx.Reply(fmt.Sprintf("Error Bot not connected"))
}
