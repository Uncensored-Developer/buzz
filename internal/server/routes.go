package server

import (
	"github.com/Uncensored-Developer/buzz/pkg/config"
	"go.uber.org/zap"
	"net/http"
)

func addRoutes(
	mux *http.ServeMux,
	cfg *config.Config,
	logger *zap.Logger,
) {
	//mux.HandleFunc("/user/create")
}
