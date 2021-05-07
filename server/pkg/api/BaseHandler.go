package api

import (
	"net/http"
	"watch2gether/pkg/datastore"
	"watch2gether/pkg/utils"
)

type BaseHandler struct {
	*datastore.Datastore
	Config *utils.Config
}

type RequestHandler = func(w http.ResponseWriter, r *http.Request) error

// The Handler struct that takes a configured Env and a function matching
// our useful signature.
type Handler struct {
	H RequestHandler
}
