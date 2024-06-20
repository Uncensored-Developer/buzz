package delivery

import (
	"context"
	"github.com/Uncensored-Developer/buzz/internal/server/dto"
	"github.com/Uncensored-Developer/buzz/internal/users/features"
	"go.uber.org/zap"
	"net/http"
	"strings"
)

func LoggedInUserOnlyMiddleware(
	ctx context.Context,
	logger *zap.Logger,
	authService *features.AuthenticationService,
	next http.Handler,
) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			authValue := r.Header.Get("Authorization")
			if authValue == "" {
				// No Authorization Header provided
				msg := "Authorization Header not found"
				dto.SendErrorJsonResponse(w, logger, msg, http.StatusUnauthorized)
				return
			}

			values := strings.Fields(authValue)
			if len(values) != 2 {
				msg := "JWT not found"
				dto.SendErrorJsonResponse(w, logger, msg, http.StatusUnauthorized)
				return
			}

			if values[0] != "Bearer" {
				msg := "Invalid Authorization Header"
				dto.SendErrorJsonResponse(w, logger, msg, http.StatusUnauthorized)
				return
			}

			user, err := authService.GetUserFromToken(ctx, values[1])
			if err != nil {
				msg := "Invalid JWT"
				dto.SendErrorJsonResponse(w, logger, msg, http.StatusUnauthorized)
				return
			}
			c := context.WithValue(r.Context(), "user", user)
			next.ServeHTTP(w, r.WithContext(c))
		},
	)
}
