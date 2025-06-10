package api

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"

	"github.com/gorilla/mux"
)

var (
	jellyfinURL = os.Getenv("JELLYFIN_URL")
	jellyfinKey = os.Getenv("JELLYFIN_KEY")
)

// http://localhost:8080/jellyfin/86a2db653e465d9efc029f47d6ffa0c0
// JellyfinProxy serves as a generic proxy for Jellyfin API, including HLS and segment files.
// Example: /jellyfin/Videos/<EpisodeId>/main.m3u8 or /jellyfin/Videos/<EpisodeId>/<segment>.ts
func (h *handler) JellyfinProxy(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	id := vars["vidoeid"]
	ext := vars["rest"]

	// Workaround: if the original request URL contains a '?' after the videoid, manually extract the query string
	// Example: /jellyfin/86a2db653e465d9efc029f47d6ffa0c0/hls1/main/1349.ts?runtimeTicks=80940000000&actualSegmentLengthTicks=38133330
	// Gorilla mux puts everything after videoid/ into ext, but strips the query string into r.URL.RawQuery
	// So, if the client uses the correct URL, r.URL.RawQuery will be set. If not, we can't recover the query params.

	if ext == "" {
		ext = "stream.mp4" // Default to main.m3u8 if no segment is specified
	}

	targetPath := fmt.Sprintf("Videos/%s/%s", id, ext)

	target, err := url.Parse(jellyfinURL)
	if err != nil {
		http.Error(w, "Invalid Jellyfin URL", http.StatusInternalServerError)
		return
	}
	target.Path += "/" + targetPath
	// Always use the original query string from the request
	target.RawQuery = r.URL.RawQuery

	log.Printf("Proxying request to Jellyfin: %s %s", r.Method, target.String())

	req, err := http.NewRequest(r.Method, target.String(), r.Body)
	if err != nil {
		http.Error(w, "Failed to create request", http.StatusInternalServerError)
		return
	}
	req.Header = r.Header.Clone()
	req.Header.Set("X-Emby-Token", jellyfinKey)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		http.Error(w, "Failed to contact Jellyfin", http.StatusBadGateway)
		return
	}
	defer resp.Body.Close()
	// For all other files (e.g., .ts), just proxy as usual
	for k, v := range resp.Header {
		for _, vv := range v {
			w.Header().Add(k, vv)
		}
	}
	w.WriteHeader(resp.StatusCode)
	io.Copy(w, resp.Body)
}
