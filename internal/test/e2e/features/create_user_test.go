package features

import (
	"context"
	"fmt"
	"github.com/Uncensored-Developer/buzz/internal/server"
	"github.com/Uncensored-Developer/buzz/internal/test/e2e"
	"github.com/Uncensored-Developer/buzz/pkg/config"
	"github.com/Uncensored-Developer/buzz/pkg/logger"
	"github.com/Uncensored-Developer/buzz/pkg/migrate"
	"github.com/Uncensored-Developer/buzz/pkg/testcontainer"
	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"
	"net/http"
	"os"
	"testing"
	"time"
)

type createUserE2eTestSuite struct {
	suite.Suite
	ctx                context.Context
	serverURL          string
	shutdownServerFunc func()
	logger             *zap.Logger
	dbContainer        *testcontainer.TestDatabase
}

// SetupSuite sets up the test suite by initializing the server and running it in the background.
// Also set up cleanup for shutting down the server
func (c *createUserE2eTestSuite) SetupSuite() {
	ctx := context.Background()
	zapLogger := logger.NewLogger()
	c.logger = zapLogger
	testDatabase, err := testcontainer.NewTestDatabase(ctx, zapLogger)
	c.Require().NoError(err)
	c.dbContainer = testDatabase

	err = migrate.Up(testDatabase.DSN, "db/migrations")
	c.Require().NoError(err)

	err = os.Setenv("DATABASE_URL", testDatabase.DSN)
	c.Require().NoError(err)

	cfg, err := config.LoadConfig()
	c.Require().NoError(err)
	ctx, cancel := context.WithCancel(ctx)
	c.shutdownServerFunc = cancel

	c.ctx = ctx
	srv, err := server.InitializeServer()
	c.Require().NoError(err)
	go srv.Run(c.ctx)

	c.serverURL = fmt.Sprintf("http://%s:%s", cfg.Host, cfg.Port)
	err = e2e.WaitForReady(c.ctx, c.logger, 5*time.Second, c.serverURL+"/health")
	if err != nil {
		c.logger.Fatal("test server failed to start", zap.Error(err))
	}
}

func (c *createUserE2eTestSuite) TearDownSuite() {
	// cancel context to shutdown server
	c.shutdownServerFunc()

	// unset env variable and shutdown test database
	err := os.Unsetenv("DATABASE_URL")
	if err != nil {
		c.logger.Fatal("failed to unset env DATABASE_URL", zap.Error(err))
	}

	err = c.dbContainer.Shutdown()
	if err != nil {
		c.logger.Fatal("failed to terminate test databased", zap.Error(err))
	}
}

func (c *createUserE2eTestSuite) TestCreateUserRouteOnlyAllowPostRequest() {
	url := fmt.Sprintf("%s/user/create", c.serverURL)

	testCases := map[string]int{
		http.MethodPut:    405,
		http.MethodDelete: 405,
		http.MethodGet:    405,
		http.MethodPost:   201,
	}
	for method, expectedStatus := range testCases {
		c.T().Run(method, func(t *testing.T) {
			req, err := http.NewRequest(method, url, nil)
			c.Require().NoError(err)

			res, err := http.DefaultClient.Do(req)
			c.Require().NoError(err)

			c.Assert().Equal(res.StatusCode, expectedStatus)
		})
	}
}

func TestCreateUserE2e(t *testing.T) {
	suite.Run(t, new(createUserE2eTestSuite))
}
