package delivery

import (
	"go.uber.org/zap"
	"net/http"
)

func HandleUserSignUp(logger *zap.Logger) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			logger.Info("Handling user signup")
		},
	)
}
