package server

import (
	"go.uber.org/zap"
	"net/http"
)

func HandleHealthCheck(logger *zap.Logger) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(200)
			_, err := w.Write([]byte("Server is Live and Ready"))
			if err != nil {
				logger.Error("Failed to write response", zap.Error(err))
			}
		},
	)
}
