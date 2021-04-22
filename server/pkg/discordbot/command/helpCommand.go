package command

import "fmt"

type HelpCmd struct{ BaseCommand }

func init() {
	Commands["help"] = &HelpCmd{BaseCommand{"This is a Help Command"}}
}
func (cmd *HelpCmd) Execute(ctx CommandCtx) error {
	msg := "The Avalible Commands are: \n"
	keys := SortKeys(Commands)
	for _, k := range keys {
		msg = msg + fmt.Sprintf("- %s : %s \n", k, Commands[k].GetHelp())
	}
	return ctx.Reply(msg)
}
