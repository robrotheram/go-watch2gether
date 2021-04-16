package main

import (
	"context"
	"flag"
	"fmt"

	"watch2gether/pkg"
	"watch2gether/pkg/api"
	discord "watch2gether/pkg/discordbot"
	"watch2gether/pkg/hub"
	"watch2gether/pkg/media"
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

	hub := hub.NewHub()
	userStore := user.NewUserStore(rethink)
	roomStore := room.NewRoomStore(rethink)
	playlistStore := media.NewPlayistStore(rethink)

	SetupDiscordBot(config, hub, roomStore, playlistStore)

	var addr = flag.String("addr", ":8080", "The addr of the  application.")
	flag.Parse() // parse the flags
	log.Println("Starting web server on", *addr)

	go hub.CleanUP(userStore)
	pkg.SetupServer(&config)

	server := api.BaseHandler{
		Hub:      hub,
		Rooms:    roomStore,
		Users:    userStore,
		Playlist: playlistStore,
	}

	if err := pkg.StartServer(*addr, userStore, &server); err != nil {
		log.Fatal("ListenAndServe:", err)
	}

}

func SetupDiscordBot(config utils.Config, hub *hub.Hub, roomStore *room.RoomStore, playist *media.PlayistStore) {

	var token = ""
	flag.StringVar(&token, "t", "", "Bot Token")
	flag.Parse()

	if config.DiscordToken != "" && token == "" {
		token = config.DiscordToken
	}

	if token != "" {
		bot, err := discord.NewDiscordBot(hub, roomStore, playist, token, config.BaseURL)
		if err != nil {
			log.Error(err)
		} else {
			bot.Start()
		}
	} else {
		log.Info("No Discord Bot token")
	}
}

// package main

// import (
// 	"encoding/json"
// 	"fmt"
// 	"io/ioutil"
// 	"net/http"

// 	"github.com/gorilla/sessions"
// 	"golang.org/x/oauth2"
// )

// var (
// 	googleOauthConfig *oauth2.Config
// 	// TODO: randomize it
// 	oauthStateString = "pseudo-random"
// )
// var Discord = oauth2.Endpoint{
// 	AuthURL:  "https://discord.com/api/oauth2/authorize",
// 	TokenURL: "https://discord.com/api/oauth2/token",
// }
// var UserURL = "https://discord.com/api/users/@me"
// var GuildsURL = UserURL + "/guilds"

// func init() {
// 	googleOauthConfig = &oauth2.Config{
// 		RedirectURL:  "http://localhost:8080/callback",

// 		Scopes:       []string{"email", "guilds"},
// 		Endpoint:     Discord,
// 	}
// }

// var store = sessions.NewCookieStore([]byte("CHANGE_ME"))

// func sessionHandler(w http.ResponseWriter, r *http.Request) {
// 	// Get a session. We're ignoring the error resulted from decoding an
// 	// existing session: Get() always returns a session, even if empty.
// 	session, _ := store.Get(r, "session-name")
// 	// Set some session values.
// 	session.Values["foo"] = "bar"
// 	session.Values[42] = 43
// 	// Save it before we write to the response/return from the handler.
// 	err := session.Save(r, w)
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}
// }

// func getsessionHandler(w http.ResponseWriter, r *http.Request) {
// 	// Get a session. We're ignoring the error resulted from decoding an
// 	// existing session: Get() always returns a session, even if empty.
// 	session, _ := store.Get(r, "session-name")
// 	// Set some session values.
// 	fmt.Fprintf(w, "Content: %s\n", session.Values)

// }

// func main() {
// 	http.HandleFunc("/", handleMain)
// 	http.HandleFunc("/session", sessionHandler)
// 	http.HandleFunc("/get-session", getsessionHandler)
// 	http.HandleFunc("/api/v1/login", handleGoogleLogin)
// 	http.HandleFunc("/callback", handleGoogleCallback)
// 	http.HandleFunc("/api/v1/user", userHandler)
// 	http.HandleFunc("/logout", logout)
// 	fmt.Println(http.ListenAndServe(":8080", nil))
// }

// func tokenToJSON(token *oauth2.Token) (string, error) {
// 	if d, err := json.Marshal(token); err != nil {
// 		return "", err
// 	} else {
// 		return string(d), nil
// 	}
// }

// func tokenFromJSON(jsonStr string) (*oauth2.Token, error) {
// 	var token oauth2.Token
// 	if err := json.Unmarshal([]byte(jsonStr), &token); err != nil {
// 		return nil, err
// 	}
// 	return &token, nil
// }

// func validateToken(r *http.Request) (*oauth2.Token, error) {
// 	session, err := store.Get(r, "session-name")
// 	if err != nil || session.Values["token"] == nil {
// 		return nil, fmt.Errorf("Invalid Session token not found")
// 	}
// 	token, err := tokenFromJSON(session.Values["token"].(string))
// 	if err != nil {
// 		return nil, fmt.Errorf("Invalid Session token could not be decoded: %w", err)
// 	}

// 	return token, nil
// }

// func handleMain(w http.ResponseWriter, r *http.Request) {
// 	var htmlIndex = `<html>
// <body>
// 	<a href="/login">Discord Log In</a>
// </body>
// </html>`

// 	fmt.Fprintf(w, htmlIndex)
// }

// func handleGoogleLogin(w http.ResponseWriter, r *http.Request) {
// 	token, err := validateToken(r)
// 	if err == nil && token != nil {
// 		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
// 		return
// 	}

// 	url := googleOauthConfig.AuthCodeURL(oauthStateString)
// 	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
// }

// func handleGoogleCallback(w http.ResponseWriter, r *http.Request) {
// 	content, err := getToken(r.FormValue("state"), r.FormValue("code"))
// 	if err != nil {
// 		fmt.Println(err.Error())
// 		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
// 		return
// 	}
// 	//Save token to session
// 	str, _ := tokenToJSON(content)
// 	session, _ := store.Get(r, "session-name")
// 	session.Values["token"] = str
// 	err = session.Save(r, w)
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}
// 	http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
// }

// func getToken(state string, code string) (*oauth2.Token, error) {
// 	if state != oauthStateString {
// 		return nil, fmt.Errorf("invalid oauth state")
// 	}

// 	token, err := googleOauthConfig.Exchange(oauth2.NoContext, code)
// 	if err != nil {
// 		return nil, fmt.Errorf("code exchange failed: %s", err.Error())
// 	}
// 	return token, nil
// }

// func userHandler(w http.ResponseWriter, r *http.Request) {
// 	token, _ := validateToken(r)
// 	user, err := getUserInfo(token.AccessToken)
// 	guilds, err := getGuildsInfo(token.AccessToken)
// 	if err != nil {
// 		fmt.Fprintf(w, "Unable to get user: %w", err)
// 		return
// 	}
// 	fmt.Fprintf(w, "User: %s\n  Guilds : %s \n", user, guilds)

// }

// func getUserInfo(accessToken string) ([]byte, error) {

// 	req, err := http.NewRequest("GET", UserURL, nil)
// 	if err != nil {
// 		return nil, err
// 	}
// 	req.Header.Set("Accept", "application/json")
// 	req.Header.Set("Authorization", "Bearer "+accessToken)
// 	client := &http.Client{}
// 	response, err := client.Do(req)

// 	if err != nil {
// 		return nil, fmt.Errorf("failed getting user info: %s", err.Error())
// 	}

// 	defer response.Body.Close()
// 	contents, err := ioutil.ReadAll(response.Body)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed reading response body: %s", err.Error())
// 	}
// 	return contents, nil
// }

// func getGuildsInfo(accessToken string) ([]byte, error) {

// 	req, err := http.NewRequest("GET", GuildsURL, nil)
// 	if err != nil {
// 		return nil, err
// 	}
// 	req.Header.Set("Accept", "application/json")
// 	req.Header.Set("Authorization", "Bearer "+accessToken)
// 	client := &http.Client{}
// 	response, err := client.Do(req)

// 	if err != nil {
// 		return nil, fmt.Errorf("failed getting user info: %s", err.Error())
// 	}

// 	defer response.Body.Close()
// 	contents, err := ioutil.ReadAll(response.Body)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed reading response body: %s", err.Error())
// 	}

// 	return contents, nil
// }

// func logout(w http.ResponseWriter, r *http.Request) {
// 	session, _ := store.Get(r, "session-name")
// 	session.Options.MaxAge = -1
// 	err := session.Save(r, w)
// 	if err != nil {
// 		fmt.Printf("failed to delete session: %v", err)
// 		return
// 	}

// 	http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
// }
