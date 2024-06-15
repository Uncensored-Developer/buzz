//go:build wireinject
// +build wireinject

package wire

import (
	"github.com/Uncensored-Developer/buzz/internal/server"
	"github.com/Uncensored-Developer/buzz/pkg/config"
	"github.com/Uncensored-Developer/buzz/pkg/logger"
	"github.com/google/wire"
)

func InitializeServer() (*server.Server, error) {
	panic(wire.Build(
		server.NewServer,
		config.LoadConfig,
		logger.NewLogger,
	))
}
