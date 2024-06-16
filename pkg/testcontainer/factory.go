package testcontainer

import (
	"context"
	"github.com/Uncensored-Developer/buzz/pkg/config"
	"github.com/Uncensored-Developer/buzz/pkg/logger"
	"github.com/pkg/errors"
)

func InitializeTestDatabase(ctx context.Context) (*TestDatabase, error) {
	cfg, err := config.LoadConfig()
	if err != nil {
		return nil, errors.Wrap(err, "load config failed")
	}
	zapLogger := logger.NewLogger(cfg)
	testDatabase, err := NewTestDatabase(ctx, zapLogger)
	if err != nil {
		return nil, err
	}
	return testDatabase, nil
}
