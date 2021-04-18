package command

import (
	"fmt"
	"watch2gether/pkg/audioBot"

	"github.com/bwmarrin/discordgo"
)

func init() {
	Commands["join"] = &JoinCmd{BaseCommand{"Join Bot to a channel"}}
}

type JoinCmd struct{ BaseCommand }

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
	err = bot.Start()
	if err != nil {
		ctx.Reply(fmt.Sprintf("Bot error %v", err))
	}
	r.RegisterBot(bot)
	ctx.Reply(fmt.Sprintf("Bot added to voice channel"))
	//_, err = ctx.Session.ChannelVoiceJoin(ctx.Guild.ID, vc, false, true)
	return nil
}
