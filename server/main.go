package main

import (
	"flag"

	"watch2gether/pkg"
	"watch2gether/pkg/api"
	"watch2gether/pkg/datastore"
	discord "watch2gether/pkg/discordbot"
	"watch2gether/pkg/utils"

	log "github.com/sirupsen/logrus"
)

func main() {

	// //media.YTPlayist("https://www.youtube.com/watch?v=lijFuLFI9RQ")
	// media.YTPlayist("https://www.youtube.com/watch?v=iyLdoQGBchQ&list=PLnfcpZm6el8iug2BBqZev2IDCt4I_sCH0")

	// return

	log.SetFormatter(&log.TextFormatter{
		//DisableColors: true,
		FullTimestamp: true,
	})

	err := utils.LoadConfig(".")
	if err != nil {
		log.Fatalf("Config Error: %v", err)
	}
	log.SetLevel(utils.GetLoglevel())

	datastore := datastore.NewDatastore(utils.Configuration)
	SetupDiscordBot(utils.Configuration, datastore)

	var addr = flag.String("addr", ":8080", "The addr of the  application.")
	flag.Parse() // parse the flags
	log.Infof("Starting web server on", *addr)

	datastore.StartCleanUP()
	pkg.SetupServer(&utils.Configuration)

	server := api.BaseHandler{
		Datastore: datastore,
		Config:    &utils.Configuration,
	}

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
			err := bot.Start()
			if err != nil {
				log.Error(err)
			}
		}

	} else {
		log.Info("No Discord Bot token")
	}
}
