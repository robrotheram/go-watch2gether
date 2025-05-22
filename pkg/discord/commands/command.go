package commands

import (
	"fmt"
	"w2g/pkg/controllers"
	"w2g/pkg/discord/components"
	"w2g/pkg/discord/session"

	"github.com/bwmarrin/discordgo"
)

type CommandCtx struct {
	Session     *discordgo.Session
	Guild       *discordgo.Guild
	Channel     *discordgo.Channel
	Member      *discordgo.Member
	Interaction *discordgo.Interaction
	Args        []string
	Controller  *controllers.Controller
	UserSession *session.UserSession
	// BaseURL string
}

func (ctx *CommandCtx) Reply(message string) *discordgo.InteractionResponse {
	return &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Flags:   discordgo.MessageFlagsEphemeral,
			Content: message,
		},
	}
}

func (ctx *CommandCtx) Replyf(format string, a ...interface{}) *discordgo.InteractionResponse {
	return &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Flags:   discordgo.MessageFlagsEphemeral,
			Content: fmt.Sprintf(format, a...),
		},
	}
}

func (ctx *CommandCtx) Errorf(format string, a ...interface{}) *discordgo.InteractionResponse {
	return ctx.Replyf(format, a...)
}

func (ctx *CommandCtx) CmdReplyEmbed(message *components.EmbededMessage) *discordgo.InteractionResponse {
	embeds := []*discordgo.MessageEmbed{}
	embeds = append(embeds, &message.MessageEmbed)
	return ctx.CmdReplyEmbeds(embeds)
}
func (ctx *CommandCtx) CmdReplyEmbeds(embeds []*discordgo.MessageEmbed) *discordgo.InteractionResponse {
	return &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Flags:  discordgo.MessageFlagsEphemeral,
			Embeds: embeds,
		},
	}
}

func (ctx *CommandCtx) CmdReplyData(data *discordgo.InteractionResponseData) *discordgo.InteractionResponse {
	return &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: data,
	}
}

type Command struct {
	Name               string
	ApplicationCommand []discordgo.ApplicationCommand
	Aliases            []string
	Usage              string
	Function           func(ctx CommandCtx) *discordgo.InteractionResponse
}

var commands = make(map[string]Command)

func register(c ...Command) {
	for _, cmd := range c {
		commands[cmd.Name] = cmd
	}
}

func GetCommands() map[string]Command {
	return commands
}

func GetCommand(name string) (Command, error) {
	for key, cmd := range commands {
		if cmd.Function == nil {
			continue
		}
		if key == name {
			return cmd, nil
		}
		for _, alias := range cmd.Aliases {
			if alias == name {
				return cmd, nil
			}
		}
	}
	return Command{}, fmt.Errorf("error: Command `%s` not found or not implemneted yet. Stay tuned", name)
}
