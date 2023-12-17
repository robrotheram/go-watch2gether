package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"w2g/pkg/controllers"
	"w2g/pkg/media"

	"github.com/gorilla/mux"
)

type handler struct {
	*controllers.Hub
}

var channelNotFound = errorMessage(http.StatusNotFound, "Unable to find not found")
var userNotFound = errorMessage(http.StatusNotFound, "Unable to find user")

func NewHandler(hub *controllers.Hub) handler {
	return handler{
		Hub: hub,
	}
}

func (h *handler) getUser(r *http.Request) (User, error) {
	ctx := r.Context().Value("user")
	if userData, ok := ctx.(User); ok {
		return userData, nil
	}
	return User{}, fmt.Errorf("unable to get user from context")
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
	controller.Add(url, user.Username)
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
		w.Header().Set("Access-Control-Allow-Origin", "*")
		socket, err := Upgrader.Upgrade(w, r, nil)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("connection is not using the websocket protocol"))
			return
		}
		controller.AddListner(NewClient(socket, controller))
	})
}
