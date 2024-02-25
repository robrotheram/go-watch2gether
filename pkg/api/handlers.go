package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
	"w2g/pkg/api/ui"
	"w2g/pkg/controllers"
	"w2g/pkg/media"
	"w2g/pkg/playlists"
	"w2g/pkg/utils"

	"github.com/gorilla/mux"
)

type handler struct {
	*controllers.Hub
}

var playlistNotFound = errorMessage(http.StatusNotFound, "Unable to find playlist")
var channelNotFound = errorMessage(http.StatusNotFound, "Unable to find channel")
var userNotFound = errorMessage(http.StatusNotFound, "Unable to find user")

func NewHandler(hub *controllers.Hub) handler {
	return handler{
		Hub: hub,
	}
}

func (h *handler) getUser(r *http.Request) (User, error) {
	ctx := r.Context().Value(UserKey)
	if userData, ok := ctx.(User); ok {
		return userData, nil
	}
	return User{}, fmt.Errorf("unable to get user from context")
}

func (h *handler) handleGetSettings(w http.ResponseWriter, r *http.Request) {
	settings := ui.Settings{
		BotId: utils.Configuration.DiscordClientID,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(settings)
}

func (h *handler) handleCreateChannel(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	controller, err := h.Get(id)
	if err == nil {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(controller.State())
		return
	}
	controller = h.New(id)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(controller.State())
}

func (h *handler) handleGetChannel(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	controller, err := h.Get(id)
	if err != nil {
		WriteError(w, channelNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(controller.State())
}

func (h *handler) handleNextVideo(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	controller, err := h.Get(id)
	if err != nil {
		WriteError(w, channelNotFound)
		return
	}
	user, err := h.getUser(r)
	if err != nil {
		WriteError(w, userNotFound)
		return
	}
	controller.Skip(user.Username)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(controller.State())
}

func (h *handler) handleShuffleVideo(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	controller, err := h.Get(id)
	if err != nil {
		WriteError(w, channelNotFound)
		return
	}
	user, err := h.getUser(r)
	if err != nil {
		WriteError(w, userNotFound)
		return
	}
	controller.Shuffle(user.Username)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(controller.State())
}

func (h *handler) handleClearVideo(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	controller, err := h.Get(id)
	if err != nil {
		WriteError(w, channelNotFound)
		return
	}
	user, err := h.getUser(r)
	if err != nil {
		WriteError(w, userNotFound)
		return
	}
	controller.UpdateQueue([]media.Media{}, user.Username)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(controller.State())
}

func (h *handler) handleLoopVideo(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	controller, err := h.Get(id)
	if err != nil {
		WriteError(w, channelNotFound)
		return
	}
	user, err := h.getUser(r)
	if err != nil {
		WriteError(w, userNotFound)
		return
	}
	controller.Loop(user.Username)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(controller.State())
}
func (h *handler) handlePlayVideo(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	controller, err := h.Get(id)
	if err != nil {
		WriteError(w, channelNotFound)
		return
	}
	user, err := h.getUser(r)
	if err != nil {
		WriteError(w, userNotFound)
		return
	}
	controller.Start(user.Username)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(controller.State())
}

func (h *handler) handleSeekVideo(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	var seconds time.Duration
	decoder := json.NewDecoder(r.Body)
	if decoder.Decode(&seconds) != nil {
		return
	}

	controller, err := h.Get(id)
	if err != nil {
		WriteError(w, channelNotFound)
		return
	}
	user, err := h.getUser(r)
	if err != nil {
		WriteError(w, userNotFound)
		return
	}
	controller.Seek(seconds, user.Username)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(controller.State())
}

func (h *handler) handlePauseVideo(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	controller, err := h.Get(id)
	if err != nil {
		WriteError(w, channelNotFound)
		return
	}
	user, err := h.getUser(r)
	if err != nil {
		WriteError(w, userNotFound)
		return
	}
	controller.Pause(user.Username)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(controller.State())
}
func (h *handler) handleAddVideo(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	var url string
	decoder := json.NewDecoder(r.Body)
	if decoder.Decode(&url) != nil {
		return
	}
	controller, err := h.Get(id)
	if err != nil {
		WriteError(w, channelNotFound)
		return
	}
	user, err := h.getUser(r)
	if err != nil {
		WriteError(w, userNotFound)
		return
	}
	controller.Add(url, false, user.Username)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(controller.State())
}

func (h *handler) handleUpateQueue(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	var videos []media.Media
	decoder := json.NewDecoder(r.Body)
	if decoder.Decode(&videos) != nil {
		return
	}
	controller, err := h.Get(id)
	if err != nil {
		WriteError(w, channelNotFound)
		return
	}
	user, err := h.getUser(r)
	if err != nil {
		WriteError(w, userNotFound)
		return
	}
	controller.UpdateQueue(videos, user.Username)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(controller.State())
}

func (h *handler) notify() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id := vars["id"]
		controller, err := h.Get(id)
		if err != nil {
			WriteError(w, channelNotFound)
			return
		}
		user, err := h.getUser(r)
		if err != nil {
			WriteError(w, userNotFound)
			return
		}
		w.Header().Set("Access-Control-Allow-Origin", "*")
		socket, err := Upgrader.Upgrade(w, r, nil)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("connection is not using the websocket protocol"))
			return
		}
		client := NewClient(socket, controller, user)
		controller.Join(client, user.Username)
		controller.AddListner(client.id, client)
	})
}

func (h *handler) handleAddFromPlaylist(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	playlistID := vars["playlist_id"]
	p, err := h.Playlists().GetById(playlistID)
	if err != nil {
		WriteError(w, playlistNotFound)
		return
	}
	controller, err := h.Get(id)
	if err != nil {
		WriteError(w, channelNotFound)
		return
	}
	user, err := h.getUser(r)
	if err != nil {
		WriteError(w, userNotFound)
		return
	}
	state := controller.State()
	state.Queue = append(state.Queue, p.Videos...)
	controller.UpdateQueue(state.Queue, user.Username)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(controller.State())
}

func (h *handler) handleGetPlaylistsByChannel(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	p, err := h.Playlists().GetByChannel(id)
	if err != nil {
		WriteError(w, playlistNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(p)
}

func (h *handler) handleGetPlaylistsByUser(w http.ResponseWriter, r *http.Request) {
	user, err := h.getUser(r)
	if err != nil {
		WriteError(w, userNotFound)
		return
	}
	p, err := h.Playlists().GetByUser(user.Username)
	if err != nil {
		WriteError(w, playlistNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(p)

}

func (h *handler) handleGetPlaylistsById(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	p, err := h.Playlists().GetById(id)
	if err != nil {
		WriteError(w, playlistNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(p)
}

func (h *handler) handleCreateNewPlaylists(w http.ResponseWriter, r *http.Request) {
	user, err := h.getUser(r)
	if err != nil {
		WriteError(w, userNotFound)
		return
	}
	var playlist playlists.Playlist
	decoder := json.NewDecoder(r.Body)
	if decoder.Decode(&playlist) != nil {
		return
	}
	playlist.User = user.Username
	h.Playlists().Create(&playlist)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(playlist)
}

func (h *handler) handleUpdatePlaylist(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	var playlist playlists.Playlist
	decoder := json.NewDecoder(r.Body)
	if decoder.Decode(&playlist) != nil {
		return
	}
	err := h.Playlists().UpdatePlaylist(id, &playlist)
	if err != nil {
		WriteError(w, playlistNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(playlist)
}

func (h *handler) handleDeletePlaylist(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	err := h.Playlists().DeletePlaylist(id)
	if err != nil {
		WriteError(w, playlistNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(id)
}

func (h *handler) handleGetPlayers(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	controller, err := h.Get(id)
	if err != nil {
		WriteError(w, playlistNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(controller.Players().GetProgress())
}

func (h *handler) handleMediaProxy(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	controller, err := h.Get(id)
	if err != nil {
		WriteError(w, channelNotFound)
		return
	}
	media := controller.State().Current
	if media.ID == "" {
		http.Error(w, "URL parameter is missing", http.StatusBadRequest)
		return
	}
	// Make an HTTP request to the provided URL
	resp, err := http.Get(media.GetAudioUrl())
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to fetch URL: %s", err), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	// Forward the response from the external URL back to the client
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to read response body: %s", err), http.StatusInternalServerError)
		return
	}

	// Set the content type based on the response header
	contentType := resp.Header.Get("Content-Type")
	if contentType != "" {
		w.Header().Set("Content-Type", contentType)
	}

	// Write the response body back to the client
	w.WriteHeader(resp.StatusCode)
	w.Write(body)
}
