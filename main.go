package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"

	"watch2gether/pkg"

	log "github.com/sirupsen/logrus"
)

// Variables used for command line parameters
var (
	Token string
)

// func init() {

// 	flag.StringVar(&Token, "t", "", "Bot Token")
// 	flag.Parse()
// 	if Token == "" {
// 		panic("Token not supplied")
// 	}
// }

func discord(hub pkg.Hub) {

	// Create a new Discord session using the provided bot token.
	dg, err := discordgo.New("Bot " + Token)
	if err != nil {
		fmt.Println("error creating Discord session,", err)
		return
	}

	// Register the messageCreate func as a callback for MessageCreate events.
	dg.AddHandler(hub.MessageCreate)

	// In this example, we only care about receiving message events.
	dg.Identify.Intents = discordgo.MakeIntent(discordgo.IntentsGuildMessages)

	// Open a websocket connection to Discord and begin listening.
	err = dg.Open()
	if err != nil {
		fmt.Println("error opening connection,", err)
		return
	}

	// Wait here until CTRL-C or other term signal is received.
	fmt.Println("Bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	// Cleanly close down the Discord session.
	dg.Close()
}

func main() {
	log.SetFormatter(&log.TextFormatter{
		DisableColors: true,
		FullTimestamp: true,
	})

	//Setup HUB
	hub := pkg.NewHub()

	var addr = flag.String("addr", ":8080", "The addr of the  application.")
	flag.Parse() // parse the flags
	log.Println("Starting web server on", *addr)

	go hub.CleanUP()
	// go func() {
	if err := pkg.StartServer(*addr, &hub); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
	// }()
	//discord(hub)
}

// !w2g play
// !w2g pause
// !w2g / !w2g status
// !w2g random
// !w2g add url
