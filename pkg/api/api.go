package api

import (
	//...

	"net/http"
	"net/url"
	"time"
	"watch2gether/pkg/api/auth"
	"watch2gether/pkg/media"
	"watch2gether/pkg/players"
	"watch2gether/pkg/playlists"
	"watch2gether/pkg/utils"

	"github.com/gorilla/sessions"
	"github.com/jellydator/ttlcache/v3"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type API struct {
	auth     *auth.DiscordAuth
	playlist *playlists.PlaylistStore
	*players.Store
	Cache *ttlcache.Cache[string, any]
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
	videos, err := media.NewVideo(url, user.Username)
	if err != nil {
		return err
	}
	controller.Queue = append(controller.Queue, videos...)
	api.Save(controller)
	return c.JSON(http.StatusOK, controller)
}

func (api *API) handleAddFromPlaylist(c echo.Context) error {
	id := c.Param("id")
	playlistID := c.Param("playlist_id")
	p, err := api.playlist.GetById(playlistID)
	if err != nil {
		return err
	}
	controller, err := api.FindChannelById(id)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "Channel not Registered")
	}
	controller.Queue = append(controller.Queue, p.Videos...)
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
	user := c.Get("user").(auth.User)
	cache := api.Cache.Get(user.ID)
	var guilds []auth.DiscordGuild
	var err error
	if cache != nil {
		guilds = cache.Value().([]auth.DiscordGuild)
	} else {
		token := c.Get("token").(string)
		guilds, err = api.auth.GetGuilds(token)
		if err != nil {
			return err
		}
		api.Cache.Set(user.ID, guilds, ttlcache.DefaultTTL)
	}
	active := []auth.DiscordGuild{}
	for _, g := range guilds {
		if _, err := api.FindChannelById(g.ID); err == nil {
			active = append(active, g)
		}
	}
	return c.JSON(200, active)
}

func (api *API) handleGetUser(c echo.Context) error {
	return c.JSON(200, c.Get("user"))
}

func (api *API) handleGetPlaylistsByUser(c echo.Context) error {
	user := c.Get("user").(auth.User)
	p, err := api.playlist.GetByUser(user.Username)
	if err != nil {
		return err
	}
	return c.JSON(200, p)

}

func (api *API) handleGetPlaylistsByChannel(c echo.Context) error {
	id := c.Param("id")
	p, err := api.playlist.GetByChannel(id)
	if err != nil {
		return err
	}
	return c.JSON(200, p)
}

func (api *API) handleGetPlaylistsById(c echo.Context) error {
	id := c.Param("id")
	p, err := api.playlist.GetById(id)
	if err != nil {
		return err
	}
	return c.JSON(200, p)
}

func (api *API) handleCreateNewPlaylists(c echo.Context) error {
	user := c.Get("user").(auth.User)
	var playlist playlists.Playlist
	if err := c.Bind(&playlist); err != nil {
		return err
	}
	playlist.User = user.Username
	api.playlist.Create(&playlist)
	return c.JSON(200, playlist)
}

func (api *API) handleUpdatePlaylist(c echo.Context) error {
	id := c.Param("id")
	var playlist playlists.Playlist
	if err := c.Bind(&playlist); err != nil {
		return err
	}
	err := api.playlist.UpdatePlaylist(id, &playlist)
	if err != nil {
		return err
	}
	return c.JSON(200, playlist)
}

func (api *API) handleDeletePlaylist(c echo.Context) error {
	id := c.Param("id")
	return api.playlist.DeletePlaylist(id)
}

func (api *API) HandleLogout(c echo.Context) error {
	user := c.Get("user").(auth.User)
	api.Cache.Delete(user.ID)
	return api.auth.HandleLogout(c)
}

func NewApi(store *players.Store, pStore *playlists.PlaylistStore) error {
	auth := auth.NewDiscordAuth(&utils.Configuration)
	cache := ttlcache.New(
		ttlcache.WithTTL[string, any](30 * time.Minute),
	)

	api := API{
		auth:     &auth,
		playlist: pStore,
		Store:    store,
		Cache:    cache,
	}
	go api.Cache.Start()

	e := echo.New()
	e.Use(session.Middleware(sessions.NewCookieStore([]byte(utils.Configuration.SessionSecret))))
	// e.Use(middleware.Logger())
	e.Use(middleware.CORS())
	e.Use(auth.Middleware())

	e.GET("/auth/login", auth.HandleLogin)
	e.GET("/auth/logout", api.HandleLogout)
	e.GET("/auth/callback", auth.HandleCallback)
	e.GET("/auth/user", api.handleGetUser)

	a := e.Group("/api")
	a.GET("/guilds", api.handleGetGuilds)
	a.GET("/channel", api.handleGetAllChannel)
	a.GET("/channel/:id", api.handleGetChannel)
	a.GET("/channel/:id/playlist", api.handleGetPlaylistsByChannel)
	a.PUT("/channel/:id/add", api.handleAddVideo)
	a.PUT("/channel/:id/add/playlist/:playlist_id", api.handleAddFromPlaylist)

	a.POST("/channel/:id/skip", api.handleNextVideo)
	a.POST("/channel/:id/shuffle", api.handleShuffleVideo)
	a.POST("/channel/:id/loop", api.handleLoopVideo)
	a.POST("/channel/:id/play", api.handlePlayVideo)
	a.POST("/channel/:id/pause", api.handlePauseVideo)
	a.POST("/channel/:id/queue", api.handleUpateQueue)

	a.GET("/playist", api.handleGetPlaylistsByUser)
	a.PUT("/playist", api.handleCreateNewPlaylists)
	a.GET("/playist/:id", api.handleGetPlaylistsById)
	a.POST("/playist/:id", api.handleUpdatePlaylist)
	a.DELETE("/playist/:id", api.handleDeletePlaylist)

	g := e.Group("*")
	if utils.Configuration.Dev {
		url1, _ := url.Parse("http://localhost:5173")
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
