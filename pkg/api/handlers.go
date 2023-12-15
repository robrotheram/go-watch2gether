package api

import (
	"encoding/json"
	"net/http"
	"w2g/pkg/controllers"

	"github.com/gorilla/mux"
)

type handler struct {
	*controllers.Hub
}

func NewHandler(hub *controllers.Hub) handler {
	return handler{
		Hub: hub,
	}
}

func (h *handler) GetState(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	controller, _ := h.Get(id)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(controller.State())
}
