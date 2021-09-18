package pkg

import (
	"encoding/json"
	"net/http"
	"watch2gether/pkg/api"

	"watch2gether/pkg/utils"

	"github.com/gorilla/handlers"
	log "github.com/sirupsen/logrus"
)

var config *utils.Config

func SetupServer(c *utils.Config) {
	config = c
}

func StartServer(connection string, hndlr *api.BaseHandler) error {
	auth := api.NewDiscordAuth(config, hndlr.Users)
	r := api.NewRouter(&auth)

	r.HandleFunc("/auth/login", auth.HandleLogin).Methods("GET")
	r.HandleFunc("/auth/logout", auth.HandleLogout).Methods("GET")
	r.HandleFunc("/auth/callback", auth.HandleCallback).Methods("GET")
	r.Register("/auth/user", "GET", true, auth.HandleUser)

	r.Register("/api/v1/room/join", "POST", true, hndlr.JoinRoom)
	r.Register("/api/v1/room/leave", "POST", true, hndlr.LeaveRoom)

	r.Register("/api/v1/room/{room_id}/playlist/{id}/load", "GET", false, hndlr.LoadPlaylist)
	r.Register("/api/v1/room/{room_id}/playlist/{id}", "POST", false, hndlr.UpdatePlaylist)
	r.Register("/api/v1/room/{room_id}/playlist/{id}", "DELETE", false, hndlr.DeletePlaylist)
	r.Register("/api/v1/room/{room_id}/playlist/{id}", "GET", false, hndlr.GetRoomPlaylist)
	r.Register("/api/v1/room/{id}/playlist", "PUT", false, hndlr.CretePlaylist)
	r.Register("/api/v1/room/{id}/playlist", "GET", false, hndlr.GetAllRoomPlaylists)

	r.Register("/api/v1/room/{id}/videos", "POST", true, hndlr.AddVideo)

	r.Handle("/api/v1/room/{id}/ws", hndlr.ConnectRoom())
	r.Register("/api/v1/room/{id}", "GET", false, hndlr.GetRoomMeta)
	r.Register("/api/v1/room/{id}", "POST", false, hndlr.UpdateRoomMeta)
	r.Register("/api/v1/room/{id}", "DELETE", false, hndlr.DeleteRoom)

	r.Register("/api/v1/status/{roomName}", "GET", false, hndlr.GetRoomMeta)
	r.Register("/api/v1/status", "GET", false, hndlr.HubStatus)

	r.Register("/config", "GET", false, hndlr.GetConfig)

	if config.Dev {
		log.Info("Starting in dev mode")
		dp := newProxy()
		r.PathPrefix("/").Handler(dp)

	} else {
		spa := spaHandler{staticPath: "ui/build", indexPath: "index.html"}
		r.PathPrefix("/").Handler(spa)
	}

	headers := handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"})
	methods := handlers.AllowedMethods([]string{"GET", "POST", "DELETE", "PUT", "HEAD", "OPTIONS"})
	origins := handlers.AllowedOrigins([]string{"*"})

	return http.ListenAndServe(connection, handlers.CORS(headers, methods, origins)(r))
}

func getPageInfo(w http.ResponseWriter, r *http.Request) error {
	vars := r.URL.Query()
	url := vars["url"][0]

	s, err := utils.Scrape(url, 1)
	if err != nil {
		return err
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(s.Preview)

	return nil
}
