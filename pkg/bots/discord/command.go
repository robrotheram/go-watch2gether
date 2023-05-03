package discordbot

import (
	"fmt"
	"sort"

	"github.com/bwmarrin/discordgo"
)

var Commands = NewCommandHelper()

type CommandHelper struct {
	Cmds map[string]CMD
}

func NewCommandHelper() *CommandHelper {
	return &CommandHelper{
		Cmds: make(map[string]CMD),
	}
}

func (cmdh *CommandHelper) SortKeys() []string {
	keys := make([]string, 0, len(cmdh.Cmds))
	for k := range cmdh.Cmds {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

func (cmdh *CommandHelper) GetCommandByID(id string) (CMD, error) {
	for _, cmd := range cmdh.Cmds {
		if cmd.ID == id {
			return cmd, nil
		}
	}
	return CMD{}, fmt.Errorf("error: Commands not found or not implemented yet. Stay tuned")
}

func (cmdh *CommandHelper) GetCommand(name string) (CMD, error) {

	for key, cmd := range cmdh.Cmds {
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
	return CMD{}, fmt.Errorf("error: Command `%s` not found or not implemneted yet. Stay tuned", name)
}

func (cmdh *CommandHelper) Register(c ...CMD) {
	for _, cmd := range c {
		cmd.ApplicationCommand.ID = fmt.Sprintf("wtg_%s", cmd.Name)
		cmdh.Cmds[cmd.Name] = cmd
	}
}

type CMD struct {
	discordgo.ApplicationCommand
	Aliases  []string
	Usage    string
	Function func(ctx CommandCtx) *discordgo.InteractionResponse
}

func (c *CMD) Format() string {
	if c.Usage == "" {
		return fmt.Sprintf(`
		**!%s**
		%s
		`, c.Name, c.Description)
	}
	return fmt.Sprintf(`
	**!%s**
	%s
	usage: `+"`%s`"+` 
	`, c.Name, c.Description, c.Usage)
}

type EmbededMessage struct {
	discordgo.MessageEmbed
}

func EmbedBuilder(title string) *EmbededMessage {
	embed := EmbededMessage{}
	embed.Title = title
	embed.Thumbnail = &discordgo.MessageEmbedThumbnail{
		URL:    "https://watch2gether.exceptionerror.io/static/media/logo.14caa8da.jpg",
		Width:  10,
		Height: 10,
	}
	embed.Color = 0x4286f4
	return &embed
}

func (embed *EmbededMessage) AddField(field discordgo.MessageEmbedField) {
	if len(field.Value) > 1000 {
		field.Value = field.Value[:1000] + "\n\n..."
	}
	embed.Fields = append(embed.Fields, &field)
}
