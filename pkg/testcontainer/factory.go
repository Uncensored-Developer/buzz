package testcontainer

import (
	"context"
	"github.com/Uncensored-Developer/buzz/pkg/logger"
)

func InitializeTestDatabase(ctx context.Context) (*TestDatabase, error) {
	zapLogger := logger.NewLogger()
	testDatabase, err := NewTestDatabase(ctx, zapLogger)
	if err != nil {
		return nil, err
	}
	return testDatabase, nil
}
