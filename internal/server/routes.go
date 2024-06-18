package server

import (
	"context"
	delivery2 "github.com/Uncensored-Developer/buzz/internal/matches/delivery"
	features2 "github.com/Uncensored-Developer/buzz/internal/matches/features"
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
	matchService *features2.MatchService,
) {
	mux.Handle("/user/create", delivery.HandleUserCreate(ctx, logger, cfg, authService))
	mux.Handle("/login", delivery.HandleUserLogin(ctx, logger, authService))
	mux.Handle("/swipe", delivery.LoggedInOnly(
		ctx, logger, authService,
		delivery2.HandleUserSwipe(ctx, logger, matchService)))
	mux.Handle("/health", HandleHealthCheck(logger))
}