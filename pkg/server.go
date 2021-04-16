package pkg

import (
	"encoding/json"
	"fmt"
	"net/http"
	"watch2gether/pkg/api"

	user "watch2gether/pkg/user"
	"watch2gether/pkg/utils"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

var config *utils.Config

func SetupServer(c *utils.Config) {
	config = c
}

func StartServer(connection string, userDB *user.UserStore, hndlr *api.BaseHandler) error {
	auth := user.NewDiscordAuth(config, userDB)

	r := mux.NewRouter()

	r.HandleFunc("/auth/login", auth.HandleLogin).Methods("GET")
	r.HandleFunc("/auth/logout", auth.HandleLogout).Methods("GET")
	r.HandleFunc("/auth/callback", auth.HandleCallback).Methods("GET")
	r.Handle("/auth/user", auth.Middleware(auth.HandleUser)).Methods("GET")

	r.Handle("/api/v1/room/join", api.Handler{hndlr.JoinRoom}).Methods("POST")
	r.Handle("/api/v1/room/leave", api.Handler{hndlr.LeaveRoom}).Methods("POST")

	r.Handle("/api/v1/room/{room_id}/playlist/{id}/load", api.Handler{hndlr.LoadPlaylist}).Methods("GET")
	r.Handle("/api/v1/room/{room_id}/playlist/{id}", api.Handler{hndlr.UpdatePlaylist}).Methods("POST")
	r.Handle("/api/v1/room/{room_id}/playlist/{id}", api.Handler{hndlr.DeletePlaylist}).Methods("DELETE")
	r.Handle("/api/v1/room/{room_id}/playlist/{id}", api.Handler{hndlr.GetRoomPlaylist}).Methods("GET")
	r.Handle("/api/v1/room/{id}/playlist", api.Handler{hndlr.CretePlaylist}).Methods("PUT")
	r.Handle("/api/v1/room/{id}/playlist", api.Handler{hndlr.GetAllRoomPlaylists}).Methods("GET")

	r.Handle("/api/v1/room/{id}/ws", hndlr.ConnectRoom())
	r.Handle("/api/v1/room/{id}", api.Handler{hndlr.GetRoomMeta}).Methods("GET")
	r.Handle("/api/v1/room/{id}", api.Handler{hndlr.UpdateRoomMeta}).Methods("POST")
	r.Handle("/api/v1/room/{id}", api.Handler{hndlr.DeleteRoom}).Methods("DELETE")

	r.Handle("/api/v1/scrape", api.Handler{getPageInfo}).Methods("GET")

	r.Handle("/api/v1/status/{roomName}", api.Handler{hndlr.GetRoomMeta})
	r.HandleFunc("/api/v1/status", hndlr.HubStatus)

	if config.Dev {
		fmt.Println("Starting in dev mode")
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
