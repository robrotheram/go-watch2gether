package api

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"watch2gether/pkg/user"
	"watch2gether/pkg/utils"

	"github.com/gorilla/sessions"
	"golang.org/x/oauth2"
)

type DiscordAuth struct {
	oauthConfig      *oauth2.Config
	oauthStateString string
	store            *sessions.CookieStore
	UserDB           *user.UserStore
}

type DiscordGuild struct {
	ID             string        `json:"id"`
	Name           string        `json:"name"`
	Icon           string        `json:"icon"`
	Owner          bool          `json:"owner"`
	Permissions    int           `json:"permissions"`
	Features       []interface{} `json:"features"`
	PermissionsNew string        `json:"permissions_new"`
}

var discordEndpoint = oauth2.Endpoint{
	AuthURL:  "https://discord.com/api/oauth2/authorize",
	TokenURL: "https://discord.com/api/oauth2/token",
}
var sessionName = "discord"

func NewDiscordAuth(conf *utils.Config, userstore *user.UserStore) DiscordAuth {
	da := DiscordAuth{}
	da.oauthConfig = &oauth2.Config{
		RedirectURL:  conf.BaseURL + "/auth/callback",
		ClientID:     conf.DiscordClientID,
		ClientSecret: conf.DiscordClientSecret,
		Scopes:       []string{"email", "guilds"},
		Endpoint:     discordEndpoint,
	}
	da.oauthStateString = utils.RandStringRunes(20)
	da.store = sessions.NewCookieStore([]byte("CHANGE_ME"))
	da.UserDB = userstore
	return da
}

func (da *DiscordAuth) validateToken(r *http.Request) (*oauth2.Token, error) {
	session, err := da.store.Get(r, sessionName)
	if err != nil || session.Values["token"] == nil {
		return nil, fmt.Errorf("Invalid Session token not found")
	}
	token, err := tokenFromJSON(session.Values["token"].(string))
	if err != nil {
		return nil, fmt.Errorf("Invalid Session token could not be decoded: %w", err)
	}

	return token, nil
}

func (da *DiscordAuth) getToken(state string, code string) (*oauth2.Token, error) {
	if state != da.oauthStateString {
		return nil, fmt.Errorf("invalid oauth state")
	}

	token, err := da.oauthConfig.Exchange(oauth2.NoContext, code)
	if err != nil {
		return nil, fmt.Errorf("code exchange failed: %s", err.Error())
	}
	return token, nil
}

func (da *DiscordAuth) Middleware(next Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token, err := da.validateToken(r)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte(err.Error()))
			return
		}
		if token == nil {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("Missing Authorization"))
			return
		}

		user, err := da.getUser(token.AccessToken)
		ctx := r.Context()
		ctx = context.WithValue(ctx, "user", user)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}

func (da *DiscordAuth) HandleLogin(w http.ResponseWriter, r *http.Request) {
	token, err := da.validateToken(r)
	if err == nil && token != nil {
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}
	url := da.oauthConfig.AuthCodeURL(da.oauthStateString)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func (da *DiscordAuth) HandleLogout(w http.ResponseWriter, r *http.Request) {
	session, err := da.store.Get(r, sessionName)
	if err != nil {
		fmt.Printf("failed to get session: %v", err)
		return
	}
	session.Values["token"] = ""
	session.Options.MaxAge = -1
	err = session.Save(r, w)
	if err != nil {
		fmt.Printf("failed to delete session: %v", err)
		return
	}
	http.Redirect(w, r, "/", http.StatusFound)
}

func (da *DiscordAuth) HandleCallback(w http.ResponseWriter, r *http.Request) {
	content, err := da.getToken(r.FormValue("state"), r.FormValue("code"))
	if err != nil {
		fmt.Println(err.Error())
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}
	//Save token to session
	str, _ := tokenToJSON(content)
	session, _ := da.store.Get(r, sessionName)
	session.Values["token"] = str
	err = session.Save(r, w)

	user, err := da.getUser(content.AccessToken)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	_, err = da.UserDB.Find(user.ID)
	if err != nil {
		da.UserDB.Create(user)
	}
	http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
}

func (da *DiscordAuth) getUser(accessToken string) (user.User, error) {
	var usr user.User
	data, err := da.getClient("https://discord.com/api/users/@me", accessToken)
	if err != nil {
		return usr, err
	}
	if err := json.Unmarshal(data, &usr); err != nil {
		return usr, err
	}
	usr.AvatarIcon = fmt.Sprintf("https://cdn.discordapp.com/avatars/%s/%s.png", usr.ID, usr.Avatar)
	usr.Type = user.USER_TYPE_DISCORD
	return usr, nil
}

func (da *DiscordAuth) getGuilds(accessToken string) ([]DiscordGuild, error) {
	var guilds []DiscordGuild
	data, err := da.getClient("https://discord.com/api/users/@me/guilds", accessToken)
	if err != nil {
		return guilds, err
	}
	fmt.Println(string(data))
	if err := json.Unmarshal(data, &guilds); err != nil {
		return nil, err
	}
	return guilds, nil
}

func (da *DiscordAuth) getClient(url string, accessToken string) ([]byte, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", "Bearer "+accessToken)
	client := &http.Client{}
	response, err := client.Do(req)

	if err != nil {
		return nil, fmt.Errorf("failed getting user info: %s", err.Error())
	}

	defer response.Body.Close()
	contents, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("failed reading response body: %s", err.Error())
	}

	return contents, nil
}

func (da *DiscordAuth) HandleUser(w http.ResponseWriter, r *http.Request) error {
	token, err := da.validateToken(r)
	if err != nil {
		return err
	}
	duser, err := da.getUser(token.AccessToken)
	if err != nil {
		return err
	}
	guilds, err := da.getGuilds(token.AccessToken)
	if err != nil {
		return err
	}
	// //In the background register guilds to the server
	// go func() {
	// 	for _, g := range guilds {
	// 		meta := room.NewMeta(g.Name, user.NewUser(duser.Username, user.USER_TYPE_DISCORD))
	// 		meta.ID = g.ID
	// 		meta.Icon = fmt.Sprintf("https://cdn.discordapp.com/icons/%s/%s.png", g.ID, g.Icon)
	// 		meta.Type = user.USER_TYPE_DISCORD
	// 		err := da.rooms.Create(meta)
	// 		if err != nil {
	// 			fmt.Println(err)
	// 		}
	// 	}
	// }()

	resp := map[string]interface{}{
		"user":   duser,
		"guilds": guilds,
	}

	w.Header().Set("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(resp)
}

func tokenToJSON(token *oauth2.Token) (string, error) {
	if d, err := json.Marshal(token); err != nil {
		return "", err
	} else {
		return string(d), nil
	}
}

func tokenFromJSON(jsonStr string) (*oauth2.Token, error) {
	var token oauth2.Token
	if err := json.Unmarshal([]byte(jsonStr), &token); err != nil {
		return nil, err
	}
	return &token, nil
}
