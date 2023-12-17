package api

import (
	"fmt"
	"net/http"
	"w2g/pkg/controllers"
	"w2g/pkg/utils"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

type App struct {
	routes *mux.Router
	auth   *DiscordAuth
	config utils.Config
}

func NewApp(config utils.Config, hub *controllers.Hub) App {
	auth := NewDiscordAuth(config, []string{
		"/api",
	})
	router := mux.NewRouter()
	router.Use(loggingMiddleware)
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
	router.HandleFunc("/api/channel/{id}/loop", handlers.handleLoopVideo).Methods("POST")
	router.HandleFunc("/api/channel/{id}/play", handlers.handlePlayVideo).Methods("POST")
	router.HandleFunc("/api/channel/{id}/pause", handlers.handlePauseVideo).Methods("POST")
	router.HandleFunc("/api/channel/{id}/queue", handlers.handleUpateQueue).Methods("POST")
	router.HandleFunc("/api/channel/{id}/add", handlers.handleAddVideo).Methods("PUT")
	if config.Dev {
		router.PathPrefix("/").Handler(newProxy())

	} else {
		spa := spaHandler{staticPath: "ui/build", indexPath: "index.html"}
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
