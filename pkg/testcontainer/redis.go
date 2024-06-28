package testcontainer

import (
	"context"
	"github.com/pkg/errors"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/redis"
	"go.uber.org/zap"
)

type TestCacheDatabase struct {
	container testcontainers.Container
	ctx       context.Context
	DSN       string
	logger    *zap.Logger
}

func NewCacheDatabase(ctx context.Context, logger *zap.Logger) (*TestCacheDatabase, error) {
	c, err := redis.RunContainer(ctx,
		testcontainers.WithImage("docker.io/redis:7"),
		redis.WithSnapshotting(10, 1),
		redis.WithLogLevel(redis.LogLevelVerbose),
	)
	if err != nil {
		return nil, errors.Wrap(err, "failed to start redis test db")
	}
	logger.Info("REDIS test database container started successfully.")

	url, err := c.ConnectionString(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get connection url")
	}

	return &TestCacheDatabase{
		container: c,
		ctx:       ctx,
		DSN:       url,
		logger:    logger,
	}, nil
}

func (t *TestCacheDatabase) Shutdown() error {
	err := t.container.Terminate(t.ctx)
	if err == nil {
		t.logger.Info("REDIS test database container shutdown successfully.")
	}
	return errors.Wrap(err, "REDIS terminate failed")
}
