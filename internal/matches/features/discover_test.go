package features_test

import (
	"context"
	"github.com/Uncensored-Developer/buzz/internal/matches/features"
	"github.com/Uncensored-Developer/buzz/internal/users/data"
	"github.com/Uncensored-Developer/buzz/internal/users/models"
	"github.com/Uncensored-Developer/buzz/pkg/config"
	"github.com/Uncensored-Developer/buzz/pkg/db"
	"github.com/Uncensored-Developer/buzz/pkg/logger"
	"github.com/Uncensored-Developer/buzz/pkg/migrate"
	"github.com/Uncensored-Developer/buzz/pkg/testcontainer"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"
	"os"
	"testing"
	"time"
)

type discoverServiceTestSuite struct {
	suite.Suite
	ctx             context.Context
	logger          *zap.Logger
	testDatabase    *testcontainer.TestDatabase
	discoverService *features.DiscoverService
	testUsers       []models.User
}

func (d *discoverServiceTestSuite) SetupSuite() {
	d.ctx = context.Background()
	d.logger = logger.NewLogger()

	var err error
	// Start up test database container
	d.testDatabase, err = testcontainer.NewTestDatabase(d.ctx, d.logger)
	d.Require().NoError(err)

	// Run migrations on test database
	err = migrate.Up(d.testDatabase.DSN, "db/migrations")
	d.Require().NoError(err)

	// This is necessary for the config.LoadConfig to pick this from the environment variables
	err = os.Setenv("DATABASE_URL", d.testDatabase.DSN)
	d.Require().NoError(err)

	cfg, err := config.LoadConfig()
	d.Require().NoError(err)

	bunDb, err := db.Connect(d.testDatabase.DSN)
	d.Require().NoError(err)

	userRepo := data.NewUserRepository(bunDb)

	d.discoverService = features.NewDiscoverService(userRepo, cfg, d.logger)

	d.testUsers = []models.User{
		{
			Email:    gofakeit.Email(),
			Password: cfg.FakeUserPassword,
			Name:     gofakeit.Name(),
			Gender:   "F",
			Dob:      time.Date(1999, 3, 24, 0, 0, 0, 0, time.UTC),
		},
		{
			Email:    gofakeit.Email(),
			Password: cfg.FakeUserPassword,
			Name:     gofakeit.Name(),
			Gender:   "M",
			Dob:      time.Date(2001, 3, 24, 0, 0, 0, 0, time.UTC),
		},
		{
			Email:    gofakeit.Email(),
			Password: cfg.FakeUserPassword,
			Name:     gofakeit.Name(),
			Gender:   "F",
			Dob:      time.Date(2004, 2, 24, 0, 0, 0, 0, time.UTC),
		},
		{
			Email:    gofakeit.Email(),
			Password: cfg.FakeUserPassword,
			Name:     gofakeit.Name(),
			Gender:   "M",
			Dob:      time.Date(1995, 3, 24, 0, 0, 0, 0, time.UTC),
		},
		{
			Email:    gofakeit.Email(),
			Password: cfg.FakeUserPassword,
			Name:     gofakeit.Name(),
			Gender:   "O",
			Dob:      time.Date(1988, 3, 24, 0, 0, 0, 0, time.UTC),
		},
	}
	_, err = bunDb.NewInsert().Model(&d.testUsers).Exec(d.ctx)
	d.Require().NoError(err)
}

func (d *discoverServiceTestSuite) TearDownSuite() {
	err := os.Unsetenv("DATABASE_URL")
	if err != nil {
		d.logger.Error("unset DATABASE_URL env failed", zap.Error(err))
	}

	err = d.testDatabase.Shutdown()
	if err != nil {
		d.logger.Error("failed to terminate database container", zap.Error(err))
	}
}

func (d *discoverServiceTestSuite) TestFetchPotentialMatches_GenderFilter() {

	testCase := map[string]struct {
		filter        features.MatchFilter
		expectedErr   error
		expectedCount int
	}{
		"male only": {
			filter:        features.MatchFilter{Gender: "M"},
			expectedErr:   nil,
			expectedCount: 2,
		},
		"female only": {
			filter:        features.MatchFilter{Gender: "F"},
			expectedErr:   nil,
			expectedCount: 1,
		},
		"others only": {
			filter:        features.MatchFilter{Gender: "O"},
			expectedErr:   nil,
			expectedCount: 1,
		},
	}

	for name, tc := range testCase {
		d.T().Run(name, func(t *testing.T) {
			users, err := d.discoverService.FetchPotentialMatches(d.ctx, d.testUsers[0].ID, tc.filter)
			d.Assert().NoError(err)
			d.Assert().Equal(tc.expectedCount, len(users))
		})
	}
}

func TestDiscoverService(t *testing.T) {
	suite.Run(t, new(discoverServiceTestSuite))
}
