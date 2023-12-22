package components

import (
	"fmt"
	"w2g/pkg/controllers"
	"w2g/pkg/discord/session"

	"github.com/bwmarrin/discordgo"
)

type HandlerCtx struct {
	Session     *discordgo.Session
	Guild       *discordgo.Guild
	Channel     *discordgo.Channel
	Args        []string
	Controller  *controllers.Controller
	UserSession *session.UserSession
	User        *discordgo.User
	// BaseURL string
}

func (ctx *HandlerCtx) UpdateMessage(data *discordgo.InteractionResponseData) *discordgo.InteractionResponse {
	return &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseUpdateMessage,
		Data: data,
	}
}

func (ctx *HandlerCtx) SendMessage(data *discordgo.InteractionResponseData) *discordgo.InteractionResponse {
	return &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: data,
	}
}

type Handler struct {
	Name     string
	Function func(ctx HandlerCtx) *discordgo.InteractionResponse
}

var commands = make(map[string]Handler)

func register(cmd Handler) {
	commands[cmd.Name] = cmd
}

func GetHandler(name string) (Handler, error) {
	for key, cmd := range commands {
		if cmd.Function == nil {
			continue
		}
		if key == name {
			return cmd, nil
		}
	}
	return Handler{}, fmt.Errorf("error: Command `%s` not found or not implemneted yet. Stay tuned", name)
}
