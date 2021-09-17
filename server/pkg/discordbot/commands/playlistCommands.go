package commands

import (
	"fmt"
	"math/rand"
	"strings"
	"watch2gether/pkg/events"
	"watch2gether/pkg/user"

	"github.com/bwmarrin/discordgo"
)

func init() {

	PlaylistCommands := NewCommandHelper()
	PlaylistCommands.Register(
		CMD{
			Command:     "load",
			Description: "Loads a Playlist",
			Usage:       "!playist load <playlist name>",
			Function:    PlaylistLoadCmd,
		},
		CMD{
			Command:     "list",
			Description: "lists all playlist",
			Usage:       "!playist list",
			Function:    PlaylistListCmd,
		},
	)

	playlistDesc := ""
	for _, cmd := range PlaylistCommands.Cmds {
		cmd.Command = "playlist " + cmd.Command
		playlistDesc += cmd.Format()
	}

	Commands.Register(
		CMD{
			Command:     "playlist",
			Description: "Commands asscoiated with playlist" + playlistDesc,
			Aliases:     []string{"link"},
			Function: func(ctx CommandCtx) error {
				if len(ctx.Args) < 1 {
					return nil
				}
				cmd, err := PlaylistCommands.GetCommand(ctx.Args[0])
				if err != nil {
					return err
				}
				ctx.Args = ctx.Args[1:]
				return cmd.Function(ctx)
			},
		},

		CMD{
			Command:     "summon",
			Description: "loads the special playlist for the user. Note must have been created ahead of time",
			Usage:       "!summon @username",
			Function:    LinkCMD,
		},
	)
}

func PlaylistLoadCmd(ctx CommandCtx) error {
	r, ok := ctx.GetHubRoom()
	meta, err := ctx.GetMeta()
	if !ok || err != nil {
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

func PlaylistListCmd(ctx CommandCtx) error {
	playlists, err := ctx.Playlist.FindByField("RoomID", ctx.Guild.ID)
	if err != nil {
		return fmt.Errorf(":x2: unable to find playlists for the room")
	}

	messageStr := ""
	for i, playlist := range playlists {
		messageStr += fmt.Sprintf("`%d` %s added by `%s` \n\n", i+1, playlist.Name, playlist.Username)
	}

	msg := EmbedBuilder("Watch2Gether Playlists")
	msg.AddField(discordgo.MessageEmbedField{
		Name:  "This room has the following playlists",
		Value: messageStr,
	})
	return ctx.ReplyEmbed(msg)
}

func SummonCmd(ctx CommandCtx) error {
	r, ok := ctx.GetHubRoom()
	meta, err := ctx.GetMeta()
	if !ok && err != nil {
		return fmt.Errorf("room %s not active", ctx.Guild.ID)
	}

	playlists, err := ctx.Playlist.FindByField("RoomID", ctx.Guild.ID)
	if err != nil {
		return fmt.Errorf("Unable to find playlists for the room")
	}

	//Get Username from discord message
	usrID := ctx.Args[0]
	usrID = strings.Trim(usrID, ">")
	usrID = strings.Trim(usrID, "<@!")
	usr, err := ctx.Session.User(usrID)
	if err != nil {
		return fmt.Errorf("unable to find user: %w", err)
	}
	playlistName := "@" + usr.Username
	for _, playlist := range playlists {
		if strings.EqualFold(strings.TrimSuffix(playlist.Name, " "), playlistName) {
			video := playlist.Videos[rand.Intn(len(playlist.Videos))]
			queue := meta.Queue
			queue = append(queue, video)
			r.HandleEvent(events.Event{
				Action:  events.EVNT_UPDATE_QUEUE,
				Watcher: user.DISCORD_BOT,
				Queue:   queue,
			})
			return ctx.Reply(fmt.Sprintf("Playing playlist for: %s", ctx.Args[0]))
		}
	}
	return ctx.Reply(fmt.Sprintf("Could not find playlist to use to summon '%s'", playlistName))
}
