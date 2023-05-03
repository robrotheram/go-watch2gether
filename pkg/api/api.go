package api

import (
	//...

	"net/http"
	"net/url"
	"watch2gether/pkg/api/auth"
	"watch2gether/pkg/media"
	"watch2gether/pkg/players"
	"watch2gether/pkg/utils"

	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type API struct {
	auth *auth.DiscordAuth
	*players.Store
}

func (api *API) handleGetChannel(c echo.Context) error {
	id := c.Param("id")
	player, err := api.FindChannelById(id)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "Channel not Registered")

	}
	return c.JSON(http.StatusOK, player)
}

func (api *API) handleUpateQueue(c echo.Context) error {
	id := c.Param("id")
	var videos []media.Media
	if err := c.Bind(&videos); err != nil {
		return err
	}
	player, err := api.FindChannelById(id)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "Channel not Registered")
	}
	player.Queue = videos
	if err := api.Save(player); err != nil {
		return err
	}
	return c.JSON(http.StatusOK, player)
}

func (api *API) handleAddVideo(c echo.Context) error {
	id := c.Param("id")
	user := c.Get("user").(auth.User)
	controller, err := api.FindChannelById(id)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "Channel not Registered")
	}
	var url string
	if err := c.Bind(&url); err != nil {
		return err
	}
	videos, _ := media.NewVideo(url, user.Username)
	controller.Queue = append(controller.Queue, videos...)
	api.Save(controller)
	return c.JSON(http.StatusOK, controller)
}

func (api *API) handleNextVideo(c echo.Context) error {
	id := c.Param("id")
	controller, _ := api.FindControllerById(id)
	controller.Skip()
	return c.JSON(http.StatusOK, controller.GetState())
}

func (api *API) handleLoopVideo(c echo.Context) error {
	id := c.Param("id")
	controller, _ := api.FindControllerById(id)
	controller.SetLoop(!controller.GetState().Loop)
	return c.JSON(http.StatusOK, controller.GetState())
}
func (api *API) handlePlayVideo(c echo.Context) error {
	id := c.Param("id")
	controller, _ := api.FindControllerById(id)
	controller.Play()
	return c.JSON(http.StatusOK, controller.GetState())
}
func (api *API) handlePauseVideo(c echo.Context) error {
	id := c.Param("id")
	controller, _ := api.FindControllerById(id)
	controller.Pause()
	return c.JSON(http.StatusOK, controller.GetState())
}

func (api *API) handleShuffleVideo(c echo.Context) error {
	id := c.Param("id")
	controller, _ := api.FindControllerById(id)
	controller.Shuffle()
	return c.JSON(http.StatusOK, controller.GetState())
}

func (api *API) handleGetAllChannel(c echo.Context) error {
	return c.JSON(http.StatusOK, api.FindAllChannels())
}

func (api *API) handleGetGuilds(c echo.Context) error {
	token := c.Get("token").(string)
	guilds, err := api.auth.GetGuilds(token)
	if err != nil {
		return err
	}
	return c.JSON(200, guilds)
}

func (api *API) handleGetUser(c echo.Context) error {
	return c.JSON(200, c.Get("user"))
}

func NewApi(store *players.Store) error {
	auth := auth.NewDiscordAuth(&utils.Configuration)

	api := API{
		auth:  &auth,
		Store: store,
	}
	e := echo.New()
	e.Use(session.Middleware(sessions.NewCookieStore([]byte("secret"))))
	e.Use(middleware.Logger())
	e.Use(middleware.CORS())
	e.Use(auth.Middleware())

	e.GET("/auth/login", auth.HandleLogin)
	e.GET("/auth/logout", auth.HandleLogout)
	e.GET("/auth/callback", auth.HandleCallback)
	e.GET("/auth/user", api.handleGetUser)

	a := e.Group("/api")
	a.GET("/guilds", api.handleGetGuilds)
	a.GET("/channel", api.handleGetAllChannel)
	a.GET("/channel/:id", api.handleGetChannel)
	a.PUT("/channel/:id/add", api.handleAddVideo)

	a.POST("/channel/:id/skip", api.handleNextVideo)
	a.POST("/channel/:id/shuffle", api.handleShuffleVideo)
	a.POST("/channel/:id/loop", api.handleLoopVideo)
	a.POST("/channel/:id/play", api.handlePlayVideo)
	a.POST("/channel/:id/pause", api.handlePauseVideo)
	a.POST("/channel/:id/queue", api.handleUpateQueue)

	g := e.Group("*")
	if utils.Configuration.Dev {
		url1, _ := url.Parse("http://localhost:5175")
		g.Use(middleware.Proxy(middleware.NewRoundRobinBalancer([]*middleware.ProxyTarget{{URL: url1}})))
	} else {
		g.Use(middleware.StaticWithConfig(middleware.StaticConfig{
			Root:       "",
			Index:      "index.html",
			HTML5:      true,
			Filesystem: http.Dir("ui/dist"),
		}))
	}
	return e.Start(":8080")
}
