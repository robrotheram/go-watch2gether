package api

import (
	"fmt"
	"net/http"
	"w2g/pkg/api/ui"
	"w2g/pkg/controllers"
	"w2g/pkg/utils"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

type App struct {
	routes *mux.Router
	config utils.Config
}

func NewApp(config utils.Config, hub *controllers.Hub) App {
	auth := NewDiscordAuth(config, []string{
		"/api",
		"/app",
	})
	router := mux.NewRouter()
	router.Use(recoveryMiddleware)
	router.Use(auth.Middleware)
	handlers := NewHandler(hub)

	router.HandleFunc("/auth/login", auth.HandleLogin).Methods("GET")
	router.HandleFunc("/auth/logout", auth.HandleLogout).Methods("GET")
	router.HandleFunc("/auth/callback", auth.HandleCallback).Methods("GET")

	router.HandleFunc("/auth/user", auth.HandleUser)
	router.HandleFunc("/api/guilds", auth.HandlerGuilds).Methods("GET")

	router.Handle("/api/channel/{id}/ws", handlers.notify())
	router.HandleFunc("/api/channel/{id}", handlers.handleGetChannel).Methods("GET")
	router.HandleFunc("/api/channel/{id}", handlers.handleCreateChannel).Methods("POST")
	router.HandleFunc("/api/channel/{id}/skip", handlers.handleNextVideo).Methods("POST")
	router.HandleFunc("/api/channel/{id}/shuffle", handlers.handleShuffleVideo).Methods("POST")
	router.HandleFunc("/api/channel/{id}/clear", handlers.handleClearVideo).Methods("POST")
	router.HandleFunc("/api/channel/{id}/loop", handlers.handleLoopVideo).Methods("POST")
	router.HandleFunc("/api/channel/{id}/play", handlers.handlePlayVideo).Methods("POST")
	router.HandleFunc("/api/channel/{id}/pause", handlers.handlePauseVideo).Methods("POST")
	router.HandleFunc("/api/channel/{id}/queue", handlers.handleUpateQueue).Methods("POST")
	router.HandleFunc("/api/channel/{id}/seek", handlers.handleSeekVideo).Methods("POST")
	router.HandleFunc("/api/channel/{id}/add", handlers.handleAddVideo).Methods("PUT")
	router.HandleFunc("/api/channel/{id}/players", handlers.handleGetPlayers).Methods("GET")

	router.HandleFunc("/api/channel/{id}/playlist", handlers.handleGetPlaylistsByChannel).Methods("GET")
	router.HandleFunc("/api/channel/{id}/add/playlist/{playlist_id}", handlers.handleAddFromPlaylist).Methods("PUT")

	router.HandleFunc("/api/playist", handlers.handleGetPlaylistsByUser).Methods("GET")
	router.HandleFunc("/api/playist", handlers.handleCreateNewPlaylists).Methods("PUT")
	router.HandleFunc("/api/playist/{id}", handlers.handleGetPlaylistsById).Methods("GET")
	router.HandleFunc("/api/playist/{id}", handlers.handleUpdatePlaylist).Methods("POST")
	router.HandleFunc("/api/playist/{id}", handlers.handleDeletePlaylist).Methods("DELETE")
	router.HandleFunc("/api/settings", handlers.handleGetSettings).Methods("GET")
	router.HandleFunc("/api/channel/{id}/proxy", handlers.handleMediaProxy).Methods("GET")

	if config.Dev {
		router.PathPrefix("/").Handler(ui.NewProxy())

	} else {
		spa := ui.SPAHandler{StaticPath: "ui/dist", IndexPath: "index.html"}
		router.PathPrefix("/").Handler(spa)
	}

	return App{
		routes: router,
		config: config,
	}
}

func (app *App) Start() error {
	headers := handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"})
	methods := handlers.AllowedMethods([]string{"GET", "POST", "DELETE", "PUT", "HEAD", "OPTIONS"})
	origins := handlers.AllowedOrigins([]string{"*"})
	addr := fmt.Sprintf(":%s", app.config.ListenPort)
	log.Infof("Starting web server on %s", addr)
	return http.ListenAndServe(addr, handlers.CORS(headers, methods, origins)(app.routes))
}
