package main

// import (
// 	"bytes"
// 	"encoding/json"
// 	"fmt"
// 	"log"
// 	"math/rand"
// 	"net/http"

// 	"github.com/gorilla/mux"
// )

// type Client struct {
// 	name   string
// 	events chan *DashBoard
// }

// type DashBoard struct {
// 	User uint
// }

// func updateDashboard(client *Client) {
// 	for {
// 		db := &DashBoard{
// 			User: uint(rand.Uint32()),
// 		}
// 		client.events <- db
// 	}
// }

// func main() {
// 	r := mux.NewRouter()
// 	r.HandleFunc("/sse", dashboardHandler)

// 	// headersOk := handlers.AllowedHeaders([]string{"X-Requested-With, Content-Type, Authorization, Access-Control-Allow-Credentials"})
// 	// originsOk := handlers.AllowedOrigins([]string{"*"})
// 	// methodsOk := handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "OPTIONS"})

// 	log.Fatal(http.ListenAndServe(":8000", (r)))
// }

// func dashboardHandler(w http.ResponseWriter, r *http.Request) {
// 	client := &Client{name: r.RemoteAddr, events: make(chan *DashBoard, 10)}

// 	go updateDashboard(client)

// 	w.Header().Set("Content-Type", "text/event-stream")
// 	w.Header().Set("Cache-Control", "no-cache")
// 	w.Header().Set("Connection", "keep-alive")
// 	w.Header().Set("Access-Control-Allow-Origin", r.Header.Get("Origin"))
// 	w.Header().Set("Access-Control-Allow-Credentials", "true")

// 	//timeout := time.After(1 * time.Second)

// 	for ev := range client.events {
// 		var buf bytes.Buffer
// 		enc := json.NewEncoder(&buf)
// 		enc.Encode(ev)
// 		fmt.Fprintf(w, "data: %v\n\n", buf.String())
// 		fmt.Printf("data: %v\n", buf.String())
// 	}

// 	if f, ok := w.(http.Flusher); ok {
// 		f.Flush()
// 	}
// }

import (
	"os"
	"path/filepath"
	"watch2gether/pkg/api"
	discordbot "watch2gether/pkg/bots/discord"
	"watch2gether/pkg/channels"
	"watch2gether/pkg/playlists"
	"watch2gether/pkg/utils"

	"github.com/bwmarrin/discordgo"
	log "github.com/sirupsen/logrus"
)

func main() {
	err := utils.LoadConfig(".")
	if err != nil {
		log.Fatalf("Config Error: %v", err)
	}

	log.Println("Using Database Path:" + utils.Configuration.DatabasePath)
	os.MkdirAll(utils.Configuration.DatabasePath, os.ModePerm)
	path := filepath.Join(utils.Configuration.DatabasePath, "watch2gether.db")

	store, err := channels.NewStore(path)
	if err != nil {
		log.Fatalf("Invalid datastore parameters: %v", err)
	}
	playlistStore := playlists.NewPlaylistStore(store.DB)

	s, err := discordgo.New("Bot " + utils.Configuration.DiscordToken)
	if err != nil {
		log.Fatalf("Invalid bot parameters: %v", err)
	}

	s.AddHandler(discordbot.RegisterCommandHandler(store, playlistStore))
	s.Identify.Intents = discordgo.IntentsGuildMessages | discordgo.IntentsGuilds | discordgo.IntentsGuildMessages | discordgo.IntentsGuildVoiceStates
	err = s.Open()
	if err != nil {
		log.Fatalf("Invalid opening session: %v", err)
	}
	go discordbot.RegisterCommands(s)
	log.Println(api.NewApi(store, playlistStore))
	log.Println("Graceful shutdown")
}
