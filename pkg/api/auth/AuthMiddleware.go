package auth

import (
	"context"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"watch2gether/pkg/utils"

	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	log "github.com/sirupsen/logrus"
	"golang.org/x/oauth2"
)

func init() {
	gob.Register([]DiscordGuild(nil))
}

type User struct {
	ID         string `storm:"id" json:"id"`
	Username   string `json:"username"`
	Type       string `json:"type"`
	Avatar     string `json:"avatar"`
	AvatarIcon string `json:"avatar_icon"`
}

type DiscordAuth struct {
	oauthConfig      *oauth2.Config
	oauthStateString string
	// store            *sessions.CookieStore
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

var discordEndpoint = oauth2.Endpoint{
	AuthURL:  "https://discord.com/api/oauth2/authorize",
	TokenURL: "https://discord.com/api/oauth2/token",
}

func NewDiscordAuth(conf *utils.Config) DiscordAuth {
	da := DiscordAuth{}
	da.oauthConfig = &oauth2.Config{
		RedirectURL:  conf.BaseURL + "/auth/callback",
		ClientID:     conf.DiscordClientID,
		ClientSecret: conf.DiscordClientSecret,
		Scopes:       []string{"identify", "guilds"},
		Endpoint:     discordEndpoint,
	}
	da.oauthStateString = utils.RandStringRunes(20)
	return da
}

func (da *DiscordAuth) validateToken(c echo.Context) (*oauth2.Token, error) {
	jsonStr := c.Get("token")
	if jsonStr == nil {
		return nil, fmt.Errorf("invalid session token could not be decoded")
	}
	token, err := tokenFromJSON(jsonStr.(string))
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

func (da *DiscordAuth) canSkip(c echo.Context) bool {
	if c.Path() == "/auth/login" || c.Path() == "/auth/callback" || c.Path() == "/*/*" {
		return true
	}
	return false
}
func (da *DiscordAuth) Middleware() echo.MiddlewareFunc {
	// return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	// 	token, err := da.validateToken(r)
	// 	if err != nil {
	// 		w.WriteHeader(http.StatusUnauthorized)
	// 		w.Write([]byte(err.Error()))
	// 		return
	// 	}
	// 	if token == nil {
	// 		w.WriteHeader(http.StatusUnauthorized)
	// 		w.Write([]byte("Missing Authorization"))
	// 		return
	// 	}

	// 	user, _ := da.getUser(token.AccessToken)
	// 	ctx := r.Context()
	// 	ctx = context.WithValue(ctx, "user", user)
	// 	r = r.WithContext(ctx)

	// 	next.ServeHTTP(w, r)
	// })

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if da.canSkip(c) {
				return next(c)
			}
			sess, _ := session.Get("session", c)
			token := sess.Values["token"]
			if token == nil {
				fmt.Println("no token")
				return echo.NewHTTPError(http.StatusUnauthorized, "unathorized")
			}
			user, _ := da.getUser(token.(string))
			c.Set("user", user)
			c.Set("token", token.(string))
			da.validateToken(c)
			return next(c)
		}
	}
}

func (da *DiscordAuth) HandleLogin(c echo.Context) error {
	token, err := da.validateToken(c)
	if err == nil && token != nil {
		return c.Redirect(http.StatusTemporaryRedirect, "/app")
	}
	url := da.oauthConfig.AuthCodeURL(da.oauthStateString)
	return c.Redirect(http.StatusTemporaryRedirect, url)
}

// func (da *DiscordAuth) ClearSession(w http.ResponseWriter, r *http.Request) {
// 	session, err := da.store.Get(r, sessionName)
// 	if err != nil {
// 		fmt.Printf("failed to get session: %v", err)
// 		return
// 	}
// 	session.Values["token"] = ""
// 	session.Options.MaxAge = -1
// 	err = session.Save(r, w)
// 	if err != nil {
// 		fmt.Printf("failed to delete session: %v", err)
// 		return
// 	}
// 	http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
// }

func (da *DiscordAuth) HandleLogout(c echo.Context) error {
	sess, _ := session.Get("session", c)
	sess.Values = nil
	sess.Save(c.Request(), c.Response())
	return c.Redirect(http.StatusTemporaryRedirect, "/")
}

func (da *DiscordAuth) HandleCallback(c echo.Context) error {
	state := c.QueryParam("state")
	code := c.QueryParam("code")

	content, err := da.getToken(state, code)
	if err != nil {
		log.Debug(err.Error())
		return c.Redirect(http.StatusTemporaryRedirect, "/")
	}
	c.Set("token", content.AccessToken)
	t := c.Get("token")
	fmt.Println(t)
	redirect := "/app"

	sess, _ := session.Get("session", c)
	sess.Values["token"] = content.AccessToken
	sess.Save(c.Request(), c.Response())

	// session, _ := da.store.Get(r, sessionName)
	// next := (session.Values["next"]).(string)
	// if len(next) > 1 {
	// 	redirect = next
	// }
	//Save token to session
	// str, _ := tokenToJSON(content)
	// session.Values["token"] = str
	// session.Save(r, w)

	_, err = da.getUser(content.AccessToken)
	if err != nil {
		// http.Error(w, err.Error(), http.StatusInternalServerError)
		return err
	}

	// _, err = da.UserDB.FindById(user.ID)
	// if err != nil {
	// 	da.UserDB.Create(user)
	// }
	return c.Redirect(http.StatusTemporaryRedirect, redirect)
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
	// usr.Type = user.USER_TYPE_DISCORD
	return usr, nil
}

func (da *DiscordAuth) GetGuilds(accessToken string) ([]DiscordGuild, error) {
	var guilds []DiscordGuild
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

func (da *DiscordAuth) HandleUser(c echo.Context) error {
	token, err := da.validateToken(c)
	if err != nil {
		return err
	}
	dUser, err := da.getUser(token.AccessToken)
	if err != nil {
		return err
	}
	guilds, err := da.GetGuilds(token.AccessToken)
	if err != nil {
		return err
	}
	resp := map[string]interface{}{
		"user":   dUser,
		"guilds": guilds,
	}

	return c.JSON(http.StatusOK, resp)
}

// func tokenToJSON(token *oauth2.Token) (string, error) {
// 	if d, err := json.Marshal(token); err != nil {
// 		return "", err
// 	} else {
// 		return string(d), nil
// 	}
// }

func tokenFromJSON(jsonStr string) (*oauth2.Token, error) {
	var token oauth2.Token
	if err := json.Unmarshal([]byte(jsonStr), &token); err != nil {
		return nil, err
	}
	return &token, nil
}

// func getUser(r *http.Request) user.User {
// 	usr := r.Context().Value("user")
// 	return usr.(user.User)
// }
