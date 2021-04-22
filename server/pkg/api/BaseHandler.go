package api

import (
	"net/http"
	"watch2gether/pkg/datastore"
)

type BaseHandler struct {
	*datastore.Datastore
}

type RequestHandler = func(w http.ResponseWriter, r *http.Request) error

// The Handler struct that takes a configured Env and a function matching
// our useful signature.
type Handler struct {
	H RequestHandler
}
