package api

import (
	"github.com/gorilla/mux"
)

type Router struct {
	*mux.Router
	auth *DiscordAuth
}

func (r *Router) Register(path string, method string, auth bool, handler RequestHandler) {
	h := Handler{handler}
	if auth {
		r.Handle(path, r.auth.Middleware(h)).Methods(method)
	} else {
		r.Handle(path, h).Methods(method)
	}
}

func NewRouter(auth *DiscordAuth) *Router {
	router := Router{Router: mux.NewRouter(), auth: auth}
	return &router
}
