package main

import (
	"os"
	"os/signal"
	"w2g/pkg/discord"
	"w2g/pkg/utils"

	log "github.com/sirupsen/logrus"
)

func main() {
	log.SetFormatter(&log.TextFormatter{
		FullTimestamp: true,
	})
	err := utils.LoadConfig(".")
	if err != nil {
		log.Fatalf("Config Error: %v", err)
	}
	log.SetLevel(utils.Configuration.GetLoglevel())

	bot, _ := discord.NewDiscordBot(utils.Configuration)
	bot.Start()
	defer bot.Close()
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	log.Println("Press Ctrl+C to exit")
	<-stop
}
