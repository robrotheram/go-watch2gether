package command

import (
	"fmt"
	"strings"
)

type PlaylistCmd struct {
	SubCommands map[string]Command
}
type PlaylistLoadCmd struct{ BaseCommand }

func init() {
	cmd := PlaylistCmd{
		SubCommands: make(map[string]Command),
	}
	cmd.SubCommands["load"] = &PlaylistLoadCmd{BaseCommand{"Loads a Playlist"}}
	Commands["playlist"] = &cmd
}

func (cmd *PlaylistCmd) GetHelp() string {
	msg := "The Avalible Commands are: \n"
	keys := SortKeys(cmd.SubCommands)
	for _, k := range keys {
		msg = msg + fmt.Sprintf("\t - %s : %s \n", k, cmd.SubCommands[k].GetHelp())
	}
	return strings.TrimSuffix(msg, "\n")
}
func (cmd *PlaylistCmd) Execute(ctx CommandCtx) error {
	if len(ctx.Args) < 1 {
		return nil
	}
	name := strings.ToLower(ctx.Args[0])
	subCmd, found := cmd.SubCommands[name]
	if !found {
		return fmt.Errorf("Command %s, not found", name)
	}
	ctx.Args = ctx.Args[1:]
	return subCmd.Execute(ctx)
}

func (cmd *PlaylistLoadCmd) Execute(ctx CommandCtx) error {
	r, ok := ctx.GetHubRoom()
	meta, err := ctx.GetMeta()
	if !ok && err != nil {
		return fmt.Errorf("room %s not active", ctx.Guild.ID)
	}
	playlists, err := ctx.Playlist.FindByField("RoomID", ctx.Guild.ID)
	if err != nil {
		return fmt.Errorf("unable to find playlists for the room")
	}

	playlistName := strings.TrimSuffix(strings.Join(ctx.Args, " "), " ")
	fmt.Printf("Searching for playlist %s: \n", playlistName)
	for _, playlist := range playlists {
		if strings.EqualFold(strings.TrimSuffix(playlist.Name, " "), playlistName) {
			meta.Queue = append(meta.Queue, playlist.Videos...)
			ctx.SaveMeta(meta)
			r.Send(meta)
			return ctx.Reply(fmt.Sprintf("Added the playlist: %s", playlistName))
		}
	}
	return ctx.Reply(fmt.Sprintf("No playlist with the name '%s' was found", playlistName))
}
