package discordbot

import (
	"fmt"
	"strings"
	"watch2gether/pkg/channels"

	"github.com/bwmarrin/discordgo"
)

func init() {
	Commands.Register(
		CMD{
			ApplicationCommand: discordgo.ApplicationCommand{
				Name:        "playlist",
				Description: "Commands associated with playlist",
				Options: []*discordgo.ApplicationCommandOption{
					{
						Name:        "load",
						Description: "load playlist",
						Options: []*discordgo.ApplicationCommandOption{
							{
								Type:        discordgo.ApplicationCommandOptionString,
								Name:        "playlist-name",
								Description: "Name of the playlist to load",
								Required:    true,
							},
						},
						Type: discordgo.ApplicationCommandOptionSubCommand,
					},
					{
						Name:        "list",
						Description: "List playlist",
						Type:        discordgo.ApplicationCommandOptionSubCommand,
					},
				},
			},
			Function: func(ctx CommandCtx) *discordgo.InteractionResponse {
				switch ctx.Args[0] {
				case "load":
					ctx.Args = ctx.Args[1:]
					return PlaylistLoadCmd(ctx)
				case "list":
					ctx.Args = ctx.Args[1:]
					return PlaylistListCmd(ctx)
				default:
					return nil
				}
			},
		},
	)
}

func PlaylistLoadCmd(ctx CommandCtx) *discordgo.InteractionResponse {
	r, err := ctx.GetChannel(ctx.Guild.ID, channels.DISCORD)
	if err != nil {
		return ctx.Errorf("Room %s not active", ctx.Guild.ID)
	}
	playlists, err := ctx.Playlists.GetByChannel(ctx.Guild.ID)
	if err != nil {
		return ctx.Errorf(":x2: unable to find playlists for the room")
	}

	playlistName := strings.TrimSuffix(strings.Join(ctx.Args, " "), " ")
	fmt.Printf("Searching for playlist %s: \n", playlistName)
	for _, playlist := range playlists {
		if strings.EqualFold(strings.TrimSuffix(playlist.Name, " "), playlistName) {
			r.Add(playlist.Videos)
			return ctx.Reply(fmt.Sprintf("Added the playlist: %s", playlistName))
		}
	}
	return ctx.Reply(fmt.Sprintf("No playlist with the name '%s' was found", playlistName))
}

func PlaylistListCmd(ctx CommandCtx) *discordgo.InteractionResponse {
	playlists, err := ctx.Playlists.GetByChannel(ctx.Guild.ID)
	if err != nil {
		return ctx.Errorf(":x2: unable to find playlists for the room")
	}

	messageStr := ""
	for i, playlist := range playlists {
		messageStr += fmt.Sprintf("`%d` %s added by `%s` \n\n", i+1, playlist.Name, playlist.User)
	}

	msg := EmbedBuilder("Watch2Gether Playlists")
	msg.AddField(discordgo.MessageEmbedField{
		Name:  "This room has the following playlists",
		Value: messageStr,
	})
	return ctx.CmdReplyEmbed(msg)
}

// func SummonCmd(ctx CommandCtx) *discordgo.InteractionResponse {
// 	r, ok := ctx.GetHubRoom()
// 	meta, err := ctx.GetMeta()
// 	if !ok && err != nil {
// 		return ctx.Errorf("room %s not active", ctx.Guild.ID)
// 	}

// 	playlists, err := ctx.Playlist.FindByRoomID(ctx.Guild.ID)
// 	if err != nil {
// 		return ctx.Errorf("Unable to find playlists for the room")
// 	}

// 	//Get Username from discord message
// 	usrID := ctx.Args[0]
// 	usrID = strings.Trim(usrID, ">")
// 	usrID = strings.Trim(usrID, "<@!")
// 	usr, err := ctx.Session.User(usrID)
// 	if err != nil {
// 		return ctx.Errorf("unable to find user: %v", err)
// 	}
// 	playlistName := "@" + usr.Username
// 	for _, playlist := range playlists {
// 		if strings.EqualFold(strings.TrimSuffix(playlist.Name, " "), playlistName) {
// 			video := playlist.Videos[rand.Intn(len(playlist.Videos))]
// 			queue := meta.Queue
// 			queue = append(queue, video)
// 			r.HandleEvent(events.Event{
// 				Action:  events.EVENT_UPDATE_QUEUE,
// 				Watcher: user.DISCORD_BOT,
// 				Queue:   queue,
// 			})
// 			return ctx.Reply(fmt.Sprintf("Playing playlist for: %s", ctx.Args[0]))
// 		}
// 	}
// 	return ctx.Reply(fmt.Sprintf("Could not find playlist to use to summon '%s'", playlistName))
// }
