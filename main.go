package main

import (
	"context"
	"flag"
	"fmt"

	"watch2gether/pkg"
	"watch2gether/pkg/room"
	"watch2gether/pkg/user"
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

	rethink, err := utils.RethinkDBConnect(config)

	if err != nil {
		log.Fatalf("DB Error: %v", err)
	}

	userStore := user.NewUserStore(rethink)
	roomStore := room.NewRoomStore(rethink)

	//Setup HUB
	hub := pkg.NewHub()

	SetupDiscordBot(config, hub, roomStore)

	var addr = flag.String("addr", ":8080", "The addr of the  application.")
	flag.Parse() // parse the flags
	log.Println("Starting web server on", *addr)

	go hub.CleanUP(userStore)

	pkg.SetupServer()

	if err := pkg.StartServer(*addr, hub, userStore, roomStore); err != nil {
		log.Fatal("ListenAndServe:", err)
	}

}

func SetupDiscordBot(config utils.Config, hub *pkg.Hub, roomStore *room.RoomStore) {

	var token = ""
	flag.StringVar(&token, "t", "", "Bot Token")
	flag.Parse()

	if config.DiscordToken != "" && token == "" {
		token = config.DiscordToken
	}

	if token != "" {
		bot, err := pkg.NewDiscordBot(hub, roomStore, token, config.BaseURL)
		if err != nil {
			log.Error(err)
		} else {
			bot.Start()
		}
	} else {
		log.Info("No Discord Bot token")
	}
}
