package command

import (
	"fmt"
	"strings"
	"watch2gether/pkg/user"
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
	for k, v := range cmd.SubCommands {
		msg = msg + fmt.Sprintf("\t - %s : %s \n", k, v.GetHelp())
	}
	return msg
}
func (cmd *PlaylistCmd) Execute(ctx CommandCtx) error {
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
	if !ok {
		return fmt.Errorf("Room %s not active", ctx.Guild.ID)
	}
	playlists, err := ctx.Playlist.FindByField("RoomID", ctx.Guild.ID)
	if err != nil {
		return fmt.Errorf("Unable to find playlists for the room")
	}

	playlistName := strings.TrimSuffix(strings.Join(ctx.Args, " "), " ")
	fmt.Printf("Searching for playlist %s: \n", playlistName)
	for _, playlist := range playlists {
		if playlist.Name == playlistName {
			queue := r.GetQueue()
			queue = append(queue, playlist.Videos...)
			r.SetQueue(queue, user.DISCORD_BOT)
			return ctx.Reply(fmt.Sprintf("Added the playlist: %s", playlistName))
		}
	}
	return ctx.Reply(fmt.Sprintf("No playlist with the name '%s'was found", playlistName))
}
