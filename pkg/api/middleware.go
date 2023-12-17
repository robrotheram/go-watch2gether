package api

import (
	"encoding/json"
	"net/http"

	log "github.com/sirupsen/logrus"
)

type ErrorResponse struct {
	Error struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	} `json:"error"`
}
type SuccessResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// func loggingMiddleware(next http.Handler) http.Handler {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		startTime := time.Now()
// 		next.ServeHTTP(w, r)
// 		duration := time.Since(startTime)
// 		log.Infof("%s %s %s %v", r.Method, r.RequestURI, r.RemoteAddr, duration)
// 	})
// }

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

func errorMessage(code int, message string) ErrorResponse {
	return ErrorResponse{
		Error: struct {
			Code    int    `json:"code"`
			Message string `json:"message"`
		}{
			Code:    http.StatusNotFound,
			Message: "Not Found",
		},
	}
}

func WriteError(w http.ResponseWriter, response ErrorResponse) {
	w.WriteHeader(http.StatusNotFound)
	json.NewEncoder(w).Encode(response)
}

func successMessage(message string) SuccessResponse {
	return SuccessResponse{
		Code:    http.StatusOK,
		Message: message,
	}
}
