package main

import (
	"os"
	"path/filepath"
	"watch2gether/pkg/api"
	discordbot "watch2gether/pkg/bots/discord"
	"watch2gether/pkg/players"
	"watch2gether/pkg/utils"

	"github.com/bwmarrin/discordgo"
	log "github.com/sirupsen/logrus"
)

func main() {
	err := utils.LoadConfig(".")
	if err != nil {
		log.Fatalf("Config Error: %v", err)
	}

	os.MkdirAll("/app", os.ModePerm)
	path := filepath.Join("", "watch2gether.db")
	store, err := players.NewStore(path)
	if err != nil {
		log.Fatalf("Invalid datastore parameters: %v", err)
	}

	s, err := discordgo.New("Bot " + utils.Configuration.DiscordToken)
	if err != nil {
		log.Fatalf("Invalid bot parameters: %v", err)
	}

	s.AddHandler(discordbot.RegisterCommandHandler(store))
	s.Identify.Intents = discordgo.IntentsGuildMessages | discordgo.IntentsGuilds | discordgo.IntentsGuildMessages | discordgo.IntentsGuildVoiceStates
	s.Open()
	go discordbot.RegisterCommands(s)
	log.Println(api.NewApi(store))
	log.Println("Graceful shutdown")
}
