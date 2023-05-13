package main

import (
	"os"
	"path/filepath"
	"watch2gether/pkg/api"
	discordbot "watch2gether/pkg/bots/discord"
	"watch2gether/pkg/channels"
	"watch2gether/pkg/playlists"
	"watch2gether/pkg/utils"

	"github.com/bwmarrin/discordgo"
	log "github.com/sirupsen/logrus"
)

func main() {
	err := utils.LoadConfig(".")
	if err != nil {
		log.Fatalf("Config Error: %v", err)
	}

	log.Println("Using Database Path:" + utils.Configuration.DatabasePath)
	os.MkdirAll(utils.Configuration.DatabasePath, os.ModePerm)
	path := filepath.Join(utils.Configuration.DatabasePath, "watch2gether.db")

	store, err := channels.NewStore(path)
	if err != nil {
		log.Fatalf("Invalid datastore parameters: %v", err)
	}
	playlistStore := playlists.NewPlaylistStore(store.DB)

	s, err := discordgo.New("Bot " + utils.Configuration.DiscordToken)
	if err != nil {
		log.Fatalf("Invalid bot parameters: %v", err)
	}

	s.AddHandler(discordbot.RegisterCommandHandler(store, playlistStore))
	s.Identify.Intents = discordgo.IntentsGuildMessages | discordgo.IntentsGuilds | discordgo.IntentsGuildMessages | discordgo.IntentsGuildVoiceStates
	s.Open()
	go discordbot.RegisterCommands(s)
	log.Println(api.NewApi(store, playlistStore))
	log.Println("Graceful shutdown")
}
