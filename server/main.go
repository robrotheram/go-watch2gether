package main

import (
	"context"
	"flag"
	"fmt"

	"watch2gether/pkg"
	"watch2gether/pkg/api"
	"watch2gether/pkg/datastore"
	discord "watch2gether/pkg/discordbot"
	"watch2gether/pkg/utils"

	"github.com/go-redis/redis/v8"
	log "github.com/sirupsen/logrus"
)

// Variables used for command line parameters

func ping(client *redis.Client) error {
	pong, err := client.Ping(context.Background()).Result()
	if err != nil {
		return err
	}
	fmt.Println(pong, err)
	// Output: PONG <nil>

	return nil
}

func main() {

	log.SetFormatter(&log.TextFormatter{
		DisableColors: true,
		FullTimestamp: true,
	})

	config, err := utils.LoadConfig(".")
	if err != nil {
		log.Fatalf("Config Error: %v", err)
	}

	datastore := datastore.NewDatastore(config)
	SetupDiscordBot(config, datastore)

	var addr = flag.String("addr", ":8080", "The addr of the  application.")
	flag.Parse() // parse the flags
	log.Println("Starting web server on", *addr)

	datastore.StartCleanUP()
	pkg.SetupServer(&config)

	server := api.BaseHandler{datastore}

	if err := pkg.StartServer(*addr, &server); err != nil {
		log.Fatal("ListenAndServe:", err)
	}

}

func SetupDiscordBot(config utils.Config, datastore *datastore.Datastore) {

	var token = ""
	flag.StringVar(&token, "t", "", "Bot Token")
	flag.Parse()

	if config.DiscordToken != "" && token == "" {
		token = config.DiscordToken
	}

	if token != "" {
		bot, err := discord.NewDiscordBot(datastore, token, config.BaseURL)
		if err != nil {
			log.Error(err)
		} else {
			bot.Start()
		}
	} else {
		log.Info("No Discord Bot token")
	}
}
