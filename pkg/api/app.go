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
	config utils.Config
}

func NewApp(config utils.Config, hub *controllers.Hub) App {
	router := mux.NewRouter()

	router.Use(loggingMiddleware)
	router.Use(recoveryMiddleware)
	router.Use(errorMiddleware)

	handlers := NewHandler(hub)

	router.HandleFunc("/{id}", handlers.GetState).Methods("GET")

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
