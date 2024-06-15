//go:build wireinject
// +build wireinject

package wire

import (
	"context"
	"github.com/Uncensored-Developer/buzz/pkg/config"
	"github.com/Uncensored-Developer/buzz/pkg/logger"
	"github.com/Uncensored-Developer/buzz/pkg/testcontainer"
	"github.com/google/wire"
)

func InitializeTestDatabase(ctx context.Context) (*testcontainer.TestDatabase, error) {
	panic(wire.Build(
		testcontainer.NewTestDatabase,
		config.LoadConfig,
		logger.NewLogger,
	))
}
