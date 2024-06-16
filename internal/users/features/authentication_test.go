package features_test

import (
	"context"
	"github.com/Uncensored-Developer/buzz/internal/users/data"
	"github.com/Uncensored-Developer/buzz/internal/users/features"
	"github.com/Uncensored-Developer/buzz/internal/users/models"
	"github.com/Uncensored-Developer/buzz/pkg/authentication"
	"github.com/Uncensored-Developer/buzz/pkg/config"
	"github.com/Uncensored-Developer/buzz/pkg/db"
	"github.com/Uncensored-Developer/buzz/pkg/hash"
	"github.com/Uncensored-Developer/buzz/pkg/logger"
	"github.com/Uncensored-Developer/buzz/pkg/migrate"
	"github.com/Uncensored-Developer/buzz/pkg/testcontainer"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/suite"
	"github.com/uptrace/bun"
	"go.uber.org/zap"
	"testing"
	"time"
)

// Be careful not to reassign these variables
var globalLogger *zap.Logger
var globalConfig *config.Config
var globalTestDatabase *testcontainer.TestDatabase
var globalDb *bun.DB

// init initializes some dependencies before running the tests
// 1. Loads the configuration using config.LoadConfig.
// 2. Creates a global logger using logger.NewLogger.
// 3. Starts a test database using testcontainer.NewTestDatabase.
// 4. Connects to the test database using db.Connect.
// 5. Runs migrations using migrate.Up.
// It is called automatically when the package is initialized.
func init() {
	ctx := context.Background()
	cfg, err := config.LoadConfig()
	if err != nil {
		globalLogger.Fatal("Failed to load config")
	}
	globalConfig = cfg
	globalLogger = logger.NewLogger(globalConfig)
	dbInstance, err := testcontainer.NewTestDatabase(ctx, globalLogger)
	if err != nil {
		globalLogger.Fatal("Failed to start test database")
	}
	globalTestDatabase = dbInstance

	bunDB, err := db.Connect(globalTestDatabase.DSN)
	if err != nil {
		globalLogger.Fatal("Failed to connect to test database")
	}
	globalDb = bunDB

	migrate.Up(globalTestDatabase.DSN, "db/migrations")
	if err != nil {
		globalLogger.Fatal("Failed to run migrations")
	}
}

const testUserEmail = "test_user@buzz.com"
const testUserPassword = "password"

type AuthenticationServiceTestSuite struct {
	suite.Suite
	ctx context.Context
}

func TestAuthenticationService(t *testing.T) {
	suite.Run(t, new(AuthenticationServiceTestSuite))
	err := globalTestDatabase.Shutdown()
	if err != nil {
		globalLogger.Error("failed to shutdown test database container")
	}
}

func (a *AuthenticationServiceTestSuite) SetupTest() {
	a.ctx = context.Background()
	user := models.User{
		Name:     "John Doe",
		Email:    testUserEmail,
		Password: testUserPassword,
		Gender:   "M",
		Dob:      time.Now(),
	}
	_, err := globalDb.NewInsert().Model(&user).Exec(a.ctx)
	if err != nil {
		globalLogger.Fatal("error saving user for test setup",
			zap.String("errMsg", err.Error()))
	}
}

func (a *AuthenticationServiceTestSuite) TearDownTest() {
	_, err := globalDb.NewDelete().Model(&models.User{}).Where(
		"email = ?", testUserEmail).Exec(a.ctx)
	if err != nil {
		globalLogger.Fatal("error deleting user for test setup")
	}
}

func (a *AuthenticationServiceTestSuite) TestSignUpWithTakenEmail() {
	authService, err := setupAuthenticationService()
	a.Require().NoError(err)

	name := "John Doe"
	gender := "M"
	dob := time.Now()
	_, err = authService.SignUp(a.ctx, dob, name, testUserEmail, testUserPassword, gender)
	a.Assert().ErrorIs(err, features.ErrEmailTaken)
}

func setupAuthenticationService() (*features.AuthenticationService, error) {
	passwordHasher := hash.NewSHA1Hasher(globalConfig.PasswordHasherSalt)
	manager, err := authentication.NewManager(globalConfig.JwtKey)
	if err != nil {
		return nil, errors.Wrap(err, "token manager error")
	}

	userRepo := data.NewUserRepository(globalDb)
	return features.NewAuthenticationService(passwordHasher, manager, userRepo, globalConfig, globalLogger), nil
}
