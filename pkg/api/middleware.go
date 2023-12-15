package api

import (
	"encoding/json"
	"net/http"
	"time"

	log "github.com/sirupsen/logrus"
)

type ErrorResponse struct {
	Error struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	} `json:"error"`
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now()
		next.ServeHTTP(w, r)
		duration := time.Since(startTime)
		log.Infof("%s %s %s %v", r.Method, r.RequestURI, r.RemoteAddr, duration)
	})
}

func recoveryMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				log.Errorf("Panic: %v", err)
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			}
		}()
		next.ServeHTTP(w, r)
	})
}

func errorMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, r)
		status := w.(interface {
			Status() int
		}).Status()

		if status >= http.StatusBadRequest {
			switch status {
			case http.StatusNotFound:
				handleNotFoundError(w, r)
			}
		}
	})
}

func handleNotFoundError(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNotFound)

	errorResponse := ErrorResponse{
		Error: struct {
			Code    int    `json:"code"`
			Message string `json:"message"`
		}{
			Code:    http.StatusNotFound,
			Message: "Not Found",
		},
	}

	json.NewEncoder(w).Encode(errorResponse)
}
