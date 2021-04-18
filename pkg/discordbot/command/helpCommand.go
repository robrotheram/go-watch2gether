package command

import "fmt"

type HelpCmd struct{ BaseCommand }

func init() {
	Commands["help"] = &HelpCmd{BaseCommand{"This is a Help Command"}}
}
func (cmd *HelpCmd) Execute(ctx CommandCtx) error {
	msg := "The Avalible Commands are: \n"
	for k, v := range Commands {
		msg = msg + fmt.Sprintf("- %s : %s \n", k, v.GetHelp())
	}
	return ctx.Reply(msg)
}
