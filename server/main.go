package main

import (
	"flag"
	"fmt"

	"watch2gether/pkg"
	"watch2gether/pkg/api"
	"watch2gether/pkg/datastore"
	discord "watch2gether/pkg/discordbot"
	"watch2gether/pkg/utils"

	log "github.com/sirupsen/logrus"
)

func main() {
	log.SetFormatter(&log.TextFormatter{
		//DisableColors: true,
		FullTimestamp: true,
	})

	err := utils.LoadConfig(".")
	if err != nil {
		log.Fatalf("Config Error: %v", err)
	}
	log.SetLevel(utils.Configuration.GetLoglevel())

	ds := datastore.NewDatastore(utils.Configuration)
	metricCollection := datastore.NewMetricCollection(*ds)
	metricCollection.Start()

	// Run Migrations in the background
	go func() { ds.RunMigrations() }()
	SetupDiscordBot(utils.Configuration, ds)

	addr := fmt.Sprintf(":%s", utils.Configuration.ListenPort)
	log.Infof("Starting web server on %s", addr)

	ds.StartCleanUP()
	pkg.SetupServer(&utils.Configuration)

	server := api.BaseHandler{
		Datastore: ds,
		Config:    &utils.Configuration,
	}

	if err := pkg.StartServer(addr, &server); err != nil {
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
		bot, err := discord.NewDiscordBot(datastore, token, config.BaseURL, config.DiscordClientID)
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
