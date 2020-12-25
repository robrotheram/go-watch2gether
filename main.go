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
var (
	Token string
)

func init() {

	flag.StringVar(&Token, "t", "", "Bot Token")
	flag.Parse()
}

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

	// //Setup HUB
	// hub := pkg.NewHub()

	// if Token != "" {
	// 	bot, err := pkg.NewDiscordBot(&hub, Token)
	// 	if err != nil {
	// 		log.Error(err)
	// 	} else {
	// 		bot.Start()
	// 	}
	// } else {
	// 	log.Info("No Discord Bot token")
	// }

	// var addr = flag.String("addr", ":8080", "The addr of the  application.")
	// flag.Parse() // parse the flags
	// log.Println("Starting web server on", *addr)

	// //go hub.CleanUP()

	// if err := pkg.StartServer(*addr, &hub, userStore); err != nil {
	// 	log.Fatal("ListenAndServe:", err)
	// }

	_, err = utils.DBConnect(config)
	rethink, err := utils.RethinkDBConnect(config)

	if err != nil {
		log.Fatalf("DB Error: %v", err)
	}

	userStore := user.NewUserStore(rethink)
	roomStore := room.NewRoomStore(rethink)

	//Setup HUB
	hub := pkg.NewHub()

	if Token != "" {
		bot, err := pkg.NewDiscordBot(&hub, roomStore, Token)
		if err != nil {
			log.Error(err)
		} else {
			bot.Start()
		}
	} else {
		log.Info("No Discord Bot token")
	}

	var addr = flag.String("addr", ":8080", "The addr of the  application.")
	flag.Parse() // parse the flags
	log.Println("Starting web server on", *addr)

	go hub.CleanUP()

	if err := pkg.StartServer(*addr, &hub, userStore, roomStore); err != nil {
		log.Fatal("ListenAndServe:", err)
	}

}
