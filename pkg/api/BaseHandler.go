package api

import (
	"net/http"
	"watch2gether/pkg/datastore"
)

type BaseHandler struct {
	*datastore.Datastore
}

// The Handler struct that takes a configured Env and a function matching
// our useful signature.
type Handler struct {
	H func(w http.ResponseWriter, r *http.Request) error
}
