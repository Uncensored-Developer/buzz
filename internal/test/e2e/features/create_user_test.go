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
	"github.com/go-resty/resty/v2"
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

// SetupSuite initializes the test suite before running tests.
// This starts the test database and runs the migrations
// Runs the server in the background
func (c *createUserE2eTestSuite) SetupSuite() {
	ctx := context.Background()
	zapLogger := logger.NewLogger()
	c.logger = zapLogger
	testDatabase, err := testcontainer.NewTestDatabase(ctx, zapLogger)
	c.Require().NoError(err)
	c.dbContainer = testDatabase

	err = migrate.Up(testDatabase.DSN, "db/migrations")
	c.Require().NoError(err)

	// This is necessary for the config.LoadConfig to pick this from the environment variables
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

	// Wait and run checks until server is ready to receive connections
	err = e2e.WaitForReady(c.ctx, c.logger, 5*time.Second, c.serverURL+"/health")
	if err != nil {
		c.logger.Fatal("test server failed to start", zap.Error(err))
	}
}

// TearDownSuite runs after running all the tests
// cancels the server context, unsets the DATABASE_URL environment variable, and shuts down the test database.
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

func (c *createUserE2eTestSuite) TestCreateRandomUser() {
	url := fmt.Sprintf("%s/user/create", c.serverURL)

	type userResponse struct {
		Id       int64  `json:"id"`
		Email    string `json:"email"`
		Password string `json:"password"`
		Name     string `json:"name"`
		Gender   string `json:"gender"`
		Age      int    `json:"age"`
	}

	// HTTP response type for user signup
	type successResponse struct {
		Result userResponse `json:"result"`
	}

	var userResp successResponse

	const expectedStatus = 201

	client := resty.New()
	res, err := client.R().SetResult(&userResp).Post(url)
	c.logger.Error("req client error", zap.Error(err))
	c.Require().NoError(err)

	c.Assert().Equal(res.StatusCode(), expectedStatus)
	c.Assert().NotEmpty(userResp)
}

func TestCreateUserE2e(t *testing.T) {
	suite.Run(t, new(createUserE2eTestSuite))
}
