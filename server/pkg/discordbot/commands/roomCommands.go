package commands

import (
	"fmt"
	"watch2gether/pkg/audioBot"
	"watch2gether/pkg/room"

	"github.com/bwmarrin/discordgo"
)

func init() {
	Register(
		CMD{
			Command:     "join",
			Description: "Summons the bot to the voice channel you are in",
			Aliases:     []string{"summon"},
			Function:    JoinCmd,
		},

		CMD{
			Command:     "leave",
			Description: "Disconnects the bot from the voice channel it is in.",
			Aliases:     []string{"disconect"},
			Function:    LeaveCmd,
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
	return "", fmt.Errorf("Channel Not found")
}

func JoinCmd(ctx CommandCtx) error {
	vc, err := GetUserVoiceChannel(ctx.Session, ctx.User.ID)
	if err != nil {
		ctx.Reply("User not connected to voice channel")
	}
	voice, err := ctx.Session.ChannelVoiceJoin(ctx.Guild.ID, vc, false, true)
	if err != nil {
		ctx.Reply("User not connected to voice channel")
	}

	r, ok := ctx.GetHubRoom()
	if !ok {
		roomMeta, err := ctx.Rooms.Find(ctx.Guild.ID)
		if err != nil {
			return ctx.Reply(fmt.Sprintf("Bot error %v", err))
		}
		roomMeta.Type = room.ROOM_TYPE_DISCORD
		if err != nil {
			ctx.Reply("Room Not found")
		}
		ctx.Rooms.Update(roomMeta)
		r = room.New(roomMeta, ctx.Rooms)
		r.PurgeUsers(true)
		ctx.Hub.AddRoom(r)
		ctx.Reply("Room has started")
	}

	bot := audioBot.NewAudioBot(vc, ctx.Channel.ID, voice, ctx.Session)
	err = r.RegisterBot(bot)
	if err != nil {
		ctx.Reply(fmt.Sprintf("Bot error %v", err))
		bot.Disconnect()
	}
	ctx.Reply(fmt.Sprintf("Bot added to the W2G room"))
	return nil
}

func LeaveCmd(ctx CommandCtx) error {
	r, ok := ctx.GetHubRoom()
	if !ok {
		return fmt.Errorf("Room %s not active", ctx.Guild.ID)
	}
	if r.Bot != nil {
		return r.Bot.Disconnect()
	}
	return ctx.Reply(fmt.Sprintf("Error Bot not connected"))
}
