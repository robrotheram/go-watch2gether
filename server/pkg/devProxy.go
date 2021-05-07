package pkg

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
)

type devProxy struct {
	proxy *httputil.ReverseProxy
}

func newProxy() devProxy {
	dp := devProxy{}
	remote, _ := url.Parse("http://localhost:3000")
	dp.proxy = httputil.NewSingleHostReverseProxy(remote)
	fmt.Printf("Starting Proxy at: %v \n", remote)
	return dp
}

func (dp devProxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Println(r.URL)
	w.Header().Set("X-Ben", "Rad")
	w.Header().Set("Cache-Control", "no-store, max-age=0")
	w.Header().Set("Cache-Control", "no-cache, private, max-age=0")
	dp.proxy.ServeHTTP(w, r)
}
