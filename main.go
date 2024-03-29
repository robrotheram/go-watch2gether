package main

import (
	"os"
	"path/filepath"
	"strings"
	"w2g/pkg/api"
	"w2g/pkg/controllers"
	"w2g/pkg/discord"
	"w2g/pkg/utils"

	log "github.com/sirupsen/logrus"
	bolt "go.etcd.io/bbolt"
)

func createStore() (*bolt.DB, error) {
	log.Println("Using Database Path:" + utils.Configuration.DatabasePath)
	os.MkdirAll(utils.Configuration.DatabasePath, os.ModePerm)
	path := filepath.Join(utils.Configuration.DatabasePath, "watch2gether.db")
	return bolt.Open(path, 0600, nil)
}

func main() {
	log.SetFormatter(&log.TextFormatter{
		FullTimestamp: true,
	})
	err := utils.LoadConfig(".")
	if err != nil {
		log.Fatalf("Config Error: %v", err)
	}
	log.SetLevel(utils.Configuration.GetLoglevel())

	db, err := createStore()
	if err != nil {
		log.Fatalf("Database Error: %v", err)
	}

	log.Infof("Version: %s", strings.TrimSuffix(utils.Version, "\r\n"))

	hub := controllers.NewHub(db)

	bot, err := discord.NewDiscordBot(utils.Configuration, hub)
	if err != nil {
		log.Warnf("discord bot error: %v", err)
	}
	err = bot.Start()
	if err != nil {
		log.Warnf("discord bot error: %v", err)
	}
	defer bot.Close()

	app := api.NewApp(utils.Configuration, hub)
	app.Start()
}
