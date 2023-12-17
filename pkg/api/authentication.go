package api

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"w2g/pkg/utils"

	"github.com/gorilla/sessions"
	log "github.com/sirupsen/logrus"
	"golang.org/x/oauth2"
)

type DiscordAuth struct {
	oauthConfig      *oauth2.Config
	oauthStateString string
	store            *sessions.CookieStore
	users            map[string]User
	paths            []string
}

type User struct {
	ID         string `storm:"id" json:"id"`
	Username   string `json:"username"`
	Type       string `json:"type"`
	Avatar     string `json:"avatar"`
	AvatarIcon string `json:"avatar_icon"`
	Token      string `json:"-"`
}

type DiscordGuild struct {
	ID             string   `json:"id"`
	Name           string   `json:"name"`
	Icon           string   `json:"icon"`
	Owner          bool     `json:"owner"`
	Permissions    int      `json:"permissions"`
	Features       []string `json:"features"`
	PermissionsNew string   `json:"permissions_new"`
}

type DiscordGuilds []struct {
	ID             string   `json:"id"`
	Name           string   `json:"name"`
	Icon           string   `json:"icon"`
	Owner          bool     `json:"owner"`
	Permissions    int      `json:"permissions"`
	Features       []string `json:"features"`
	PermissionsNew string   `json:"permissions_new"`
}

var discordEndpoint = oauth2.Endpoint{
	AuthURL:  "https://discord.com/api/oauth2/authorize",
	TokenURL: "https://discord.com/api/oauth2/token",
}
var sessionName = "discord"

func NewDiscordAuth(conf utils.Config, paths []string) DiscordAuth {
	da := DiscordAuth{}
	da.oauthConfig = &oauth2.Config{
		RedirectURL:  conf.BaseURL + "/auth/callback",
		ClientID:     conf.DiscordClientID,
		ClientSecret: conf.DiscordClientSecret,
		Scopes:       []string{"identify", "guilds"},
		Endpoint:     discordEndpoint,
	}
	da.oauthStateString = utils.RandStringRunes(20)
	da.store = sessions.NewCookieStore([]byte(conf.SessionSecret))
	da.users = make(map[string]User)
	da.paths = paths
	return da
}

func (da *DiscordAuth) validateToken(r *http.Request) (*oauth2.Token, error) {
	session, err := da.store.Get(r, sessionName)
	if err != nil || session.Values["token"] == nil {
		return nil, fmt.Errorf("invalid session token not found")
	}
	token, err := tokenFromJSON(session.Values["token"].(string))
	if err != nil {
		return nil, fmt.Errorf("invalid session token could not be decoded: %w", err)
	}

	return token, nil
}

func (da *DiscordAuth) getToken(state string, code string) (*oauth2.Token, error) {
	if state != da.oauthStateString {
		return nil, fmt.Errorf("invalid oauth state")
	}

	token, err := da.oauthConfig.Exchange(context.Background(), code)
	if err != nil {
		return nil, fmt.Errorf("code exchange failed: %s", err.Error())
	}
	return token, nil
}

func (da *DiscordAuth) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !contains(da.paths, r.URL.Path) {
			next.ServeHTTP(w, r)
			return
		}
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

		user, _ := da.getUser(token.AccessToken)
		ctx := r.Context()
		ctx = context.WithValue(ctx, "user", user)
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}

func (da *DiscordAuth) HandleLogin(w http.ResponseWriter, r *http.Request) {
	token, err := da.validateToken(r)
	next := r.URL.Query().Get("next")

	if err == nil && token != nil {
		http.Redirect(w, r, "/app", http.StatusTemporaryRedirect)
		return
	}

	session, _ := da.store.Get(r, sessionName)
	session.Values["next"] = next
	session.Save(r, w)

	url := da.oauthConfig.AuthCodeURL(da.oauthStateString)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func (da *DiscordAuth) ClearSession(w http.ResponseWriter, r *http.Request) {
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
	http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
}
func (da *DiscordAuth) HandleLogout(w http.ResponseWriter, r *http.Request) {
	da.ClearSession(w, r)
}

func (da *DiscordAuth) HandleCallback(w http.ResponseWriter, r *http.Request) {
	content, err := da.getToken(r.FormValue("state"), r.FormValue("code"))
	if err != nil {
		log.Debug(err.Error())
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	redirect := "/app"
	session, _ := da.store.Get(r, sessionName)
	next := (session.Values["next"]).(string)
	if len(next) > 1 {
		redirect = next
	}
	//Save token to session
	str, _ := tokenToJSON(content)
	session.Values["token"] = str
	session.Save(r, w)

	user, err := da.getUser(content.AccessToken)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_, err = da.FindUserById(user.ID)
	if err != nil {
		da.users[user.ID] = user
	}
	http.Redirect(w, r, redirect, http.StatusTemporaryRedirect)
}

func (da *DiscordAuth) FindUserById(id string) (User, error) {
	if _, ok := da.users[id]; !ok {
		return User{}, fmt.Errorf("User not found")
	}
	return da.users[id], nil
}

func (da *DiscordAuth) getUser(accessToken string) (User, error) {
	var usr User
	data, err := da.getClient("https://discord.com/api/users/@me", accessToken)
	if err != nil {
		return usr, err
	}
	if err := json.Unmarshal(data, &usr); err != nil {
		return usr, err
	}
	usr.AvatarIcon = fmt.Sprintf("https://cdn.discordapp.com/avatars/%s/%s.png", usr.ID, usr.Avatar)
	usr.Token = accessToken
	// usr.Type = user.USER_TYPE_DISCORD
	return usr, nil
}

func (da *DiscordAuth) getGuilds(accessToken string) (DiscordGuilds, error) {
	var guilds DiscordGuilds
	data, err := da.getClient("https://discord.com/api/users/@me/guilds", accessToken)
	if err != nil {
		return guilds, err
	}
	if err := json.Unmarshal(data, &guilds); err != nil {
		log.Info(string(data))
		log.Infof("get guilds %v", err)
		return nil, nil
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
	contents, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("failed reading response body: %s", err.Error())
	}

	return contents, nil
}

func (da *DiscordAuth) HandleUser(w http.ResponseWriter, r *http.Request) {
	token, err := da.validateToken(r)
	if err != nil {
		da.ClearSession(w, r)
		return
	}
	dUser, err := da.getUser(token.AccessToken)
	if err != nil {
		da.ClearSession(w, r)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(dUser)
}

func (da *DiscordAuth) HandlerGuilds(w http.ResponseWriter, r *http.Request) {
	token, err := da.validateToken(r)
	if err != nil {
		return
	}
	guilds, err := da.getGuilds(token.AccessToken)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(guilds)
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

func getUser(r *http.Request) User {
	usr := r.Context().Value("user")
	return usr.(User)
}

func contains(s []string, e string) bool {
	for _, a := range s {
		if strings.Contains(e, a) {
			return true
		}
	}
	return false
}
