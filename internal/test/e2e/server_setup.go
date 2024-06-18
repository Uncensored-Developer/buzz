package e2e

import (
	"context"
	"fmt"
	"github.com/Uncensored-Developer/buzz/internal/server"
	"github.com/Uncensored-Developer/buzz/pkg/config"
	"github.com/Uncensored-Developer/buzz/pkg/logger"
	"github.com/Uncensored-Developer/buzz/pkg/migrate"
	"github.com/Uncensored-Developer/buzz/pkg/testcontainer"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"os"
	"time"
)

// TestServerSuite represents utilities to ensure proper server setup for e2e tests
type TestServerSuite struct {
	Ctx                context.Context
	ServerURL          string
	ShutdownServerFunc func()
	Logger             *zap.Logger
	DbContainer        *testcontainer.TestDatabase
	Config             *config.Config
}

// StartUp setups up the app server and all it's dependencies
// Ensure this is ran before any API e2e test
func (s *TestServerSuite) StartUp() error {
	ctx := context.Background()
	zapLogger := logger.NewLogger()
	s.Logger = zapLogger

	// Start up test database container
	testDatabase, err := testcontainer.NewTestDatabase(ctx, zapLogger)
	if err != nil {
		return errors.Wrap(err, "test database failed to start")
	}
	s.DbContainer = testDatabase

	// Run migrations on test database
	err = migrate.Up(testDatabase.DSN, "db/migrations")
	if err != nil {
		return errors.Wrap(err, "migration failed")
	}

	// This is necessary for the config.LoadConfig to pick this from the environment variables
	err = os.Setenv("DATABASE_URL", testDatabase.DSN)
	if err != nil {
		return errors.Wrap(err, "set env failed")
	}

	cfg, err := config.LoadConfig()
	if err != nil {
		return errors.Wrap(err, "load config failed")
	}
	ctx, cancel := context.WithCancel(ctx)
	s.Config = cfg
	s.ShutdownServerFunc = cancel

	s.Ctx = ctx
	srv, err := server.InitializeServer()
	if err != nil {
		return errors.Wrap(err, "failed to start server")
	}
	go srv.Run(s.Ctx) // Run server in background

	// Set server url from config info
	s.ServerURL = fmt.Sprintf("http://%s:%s", cfg.Host, cfg.Port)

	// Wait and run checks until server is ready to receive connections
	err = WaitForReady(s.Ctx, s.Logger, 5*time.Second, s.ServerURL+"/health")
	if err != nil {
		return errors.Wrap(err, "test server failed to start")
	}
	return nil
}

func (s *TestServerSuite) Shutdown() error {
	// cancel context to shutdown server
	s.ShutdownServerFunc()

	// unset env variable and shutdown test database
	err := os.Unsetenv("DATABASE_URL")
	if err != nil {
		return errors.Wrap(err, "unset DATABASE_URL env failed")
	}

	err = s.DbContainer.Shutdown()
	if err != nil {
		return errors.Wrap(err, "failed to terminate database container")
	}
	return nil
}
