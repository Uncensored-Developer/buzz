package server

import (
	"context"
	"github.com/Uncensored-Developer/buzz/internal/users/delivery"
	"github.com/Uncensored-Developer/buzz/internal/users/features"
	"github.com/Uncensored-Developer/buzz/pkg/config"
	"go.uber.org/zap"
	"net/http"
)

func addRoutes(
	ctx context.Context,
	mux *http.ServeMux,
	cfg *config.Config,
	logger *zap.Logger,
	authService *features.AuthenticationService,
) {
	mux.Handle("/user/create", delivery.HandleUserCreate(ctx, logger, cfg, authService))
	mux.Handle("/health", HandleHealthCheck(logger))
}
