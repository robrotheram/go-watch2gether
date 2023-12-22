package ui

import (
	"net/http"
	"net/http/httputil"
	"net/url"

	log "github.com/sirupsen/logrus"
)

type devProxy struct {
	proxy *httputil.ReverseProxy
}

func NewProxy() devProxy {
	dp := devProxy{}
	remote, _ := url.Parse("http://localhost:5173")
	dp.proxy = httputil.NewSingleHostReverseProxy(remote)
	log.Infof("Starting Proxy at: %v", remote)
	return dp
}

func (dp devProxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Debug(r.URL)
	w.Header().Set("X-Ben", "Rad")
	w.Header().Set("Cache-Control", "no-store, max-age=0")
	w.Header().Set("Cache-Control", "no-cache, private, max-age=0")
	dp.proxy.ServeHTTP(w, r)
}
