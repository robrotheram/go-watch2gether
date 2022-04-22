package commands

import (
	"fmt"
	"strings"
	"watch2gether/pkg/audioBot"
	"watch2gether/pkg/room"
	meta "watch2gether/pkg/roomMeta"
	"watch2gether/pkg/user"

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

func CreateRoom(ctx CommandCtx) (*meta.Meta, error) {
	usr := user.NewUser(ctx.Guild.OwnerID, user.USER_TYPE_DISCORD)
	return ctx.Rooms.GetOrCreate(ctx.Guild.ID, ctx.Guild.Name, usr)
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

	r, ok := ctx.GetHubRoom()
	if !ok {
		roomMeta, err := ctx.Rooms.Find(ctx.Guild.ID)
		if err != nil {
			if strings.Contains(err.Error(), "not found") {
				roomMeta, err = CreateRoom(ctx)
				if err != nil {
					return ctx.Reply(fmt.Sprintf("Bot error: %v", err))
				}
			}
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

	bot := audioBot.NewAudioBot(vc, ctx.Channel.ID, ctx.Guild.ID, voice, ctx.Session)
	err = r.RegisterBot(bot)
	if err != nil {
		ctx.Reply(fmt.Sprintf("Bot error: %v", err))
		bot.Disconnect()
	}
	return ctx.Reply("Bot added to the room")
}

func LeaveCmd(ctx CommandCtx) *discordgo.InteractionResponse {
	r, ok := ctx.GetHubRoom()
	if !ok {
		return ctx.Reply(fmt.Sprintf("Room %s not active", ctx.Guild.ID))
	}
	if r.Bot != nil {
		r.Bot.Disconnect()
		return ctx.Reply("watch2gether has left the room")
	}
	return ctx.Reply("Error Bot not connected")
}
