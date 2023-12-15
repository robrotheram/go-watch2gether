package main

import (
	"fmt"
	"w2g/pkg/api"
	"w2g/pkg/controllers"
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
	hub := controllers.NewHub()

	bot, _ := discord.NewDiscordBot(utils.Configuration, hub)
	bot.Start()
	defer bot.Close()

	app := api.NewApp(utils.Configuration, hub)
	fmt.Println(app.Start())

	// stop := make(chan os.Signal, 1)
	// signal.Notify(stop, os.Interrupt)
	// log.Println("Press Ctrl+C to exit")
	// <-stop
}
