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
	"github.com/pkg/errors"
	"github.com/stretchr/testify/suite"
	"github.com/uber/h3-go/v4"
	"github.com/uptrace/bun"
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
	bunDb           *bun.DB
	cfg             *config.Config
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
	d.cfg = cfg

	d.bunDb, err = db.Connect(d.testDatabase.DSN)
	d.Require().NoError(err)

	userRepo := data.NewUserRepository(d.bunDb)

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
			Gender:   "M",
			Dob:      time.Date(1995, 3, 24, 0, 0, 0, 0, time.UTC),
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
			Gender:   "O",
			Dob:      time.Date(1988, 3, 24, 0, 0, 0, 0, time.UTC),
		},
	}
	_, err = d.bunDb.NewInsert().Model(&d.testUsers).Exec(d.ctx)
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

func (d *discoverServiceTestSuite) TestFetchPotentialMatches_NoRadius() {

	testCase := map[string]struct {
		filter        features.MatchFilter
		expectedErr   error
		expectedCount int
	}{
		"male only": {
			filter:        features.MatchFilter{Gender: features.MaleGender},
			expectedErr:   nil,
			expectedCount: 2,
		},
		"female only": {
			filter:        features.MatchFilter{Gender: features.FemaleGender},
			expectedErr:   nil,
			expectedCount: 1,
		},
		"others only": {
			filter:        features.MatchFilter{Gender: features.OtherGender},
			expectedErr:   nil,
			expectedCount: 1,
		},
		"minimum age only": {
			filter:        features.MatchFilter{MinAge: 25},
			expectedErr:   nil,
			expectedCount: 2,
		},
		"maximum age only": {
			filter:        features.MatchFilter{MaxAge: 24},
			expectedErr:   nil,
			expectedCount: 2,
		},
		"age range only": {
			filter:        features.MatchFilter{MinAge: 21, MaxAge: 28},
			expectedErr:   nil,
			expectedCount: 1,
		},
		"age range and gender": {
			filter:        features.MatchFilter{MinAge: 21, MaxAge: 28, Gender: features.OtherGender},
			expectedErr:   nil,
			expectedCount: 0,
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

func (d *discoverServiceTestSuite) TestFetchPotentialMatches_WithRadius() {
	testUsers, teardown, err := setupUsers(d.ctx, d.cfg, d.bunDb, d.logger)
	d.Require().NoError(err)
	defer teardown()

	filters := features.MatchFilter{
		Radius: 100,
	}
	wantCount := 3
	wantEmails := []string{"user2@buzz.com", "user3@buzz.com", "user4@buzz.com"}
	users, err := d.discoverService.FetchPotentialMatches(d.ctx, testUsers[0].ID, filters)

	var gotEmails []string
	for _, user := range users {
		gotEmails = append(gotEmails, user.Email)
	}
	d.Require().NoError(err)
	d.Assert().Equal(wantCount, len(users))
	d.Assert().ElementsMatch(wantEmails, gotEmails)
}

func TestDiscoverService(t *testing.T) {
	suite.Run(t, new(discoverServiceTestSuite))
}

func setupUsers(
	ctx context.Context,
	cfg *config.Config,
	db *bun.DB,
	logger *zap.Logger,
) ([]models.User, func() error, error) {
	testUsers := []models.User{
		{
			Email:     "user1@buzz.com",
			Password:  "fakePass",
			Name:      gofakeit.Name(),
			Gender:    "F",
			Dob:       time.Date(1999, 3, 24, 0, 0, 0, 0, time.UTC),
			Longitude: 0.5026768,
			Latitude:  51.2725887,
		},
		{
			Email:     "user2@buzz.com",
			Password:  "cfg.FakeUserPassword",
			Name:      gofakeit.Name(),
			Gender:    "M",
			Dob:       time.Date(2001, 3, 24, 0, 0, 0, 0, time.UTC),
			Latitude:  50.96284649,
			Longitude: -0.12981616,
		},
		{
			Email:     "user3@buzz.com",
			Password:  "cfg.FakeUserPassword",
			Name:      gofakeit.Name(),
			Gender:    "M",
			Dob:       time.Date(1995, 3, 24, 0, 0, 0, 0, time.UTC),
			Latitude:  51.03052403,
			Longitude: 0.18169958,
		},
		{
			Email:     "user4@buzz.com",
			Password:  "user1@buzz.com",
			Name:      gofakeit.Name(),
			Gender:    "F",
			Dob:       time.Date(2004, 2, 24, 0, 0, 0, 0, time.UTC),
			Longitude: 0.62093072,
			Latitude:  51.31488722,
		},
		{
			Email:     "user5@buzz.com",
			Password:  "cfg.FakeUserPassword",
			Name:      gofakeit.Name(),
			Gender:    "O",
			Dob:       time.Date(1988, 3, 24, 0, 0, 0, 0, time.UTC),
			Latitude:  6.524379,
			Longitude: 3.379206,
		},
	}
	_, err := db.NewInsert().Model(&testUsers).Exec(ctx)
	if err != nil {
		return nil, nil, errors.Wrap(err, "failed to create test location users")
	}

	var emails []string
	for _, user := range testUsers {
		emails = append(emails, user.Email)

		latLng := h3.NewLatLng(user.Latitude, user.Latitude)
		cell := h3.LatLngToCell(latLng, cfg.H3Resolution)
		user.H3Index = int64(cell)

		_, err = db.NewUpdate().Model(&user).WherePK().Returning("*").Exec(ctx)
		if err != nil {
			return nil, nil, errors.Wrap(err, "failed to update h3_index")
		}
	}
	return testUsers, func() error {
		_, err := db.NewDelete().Model(&models.User{}).Where("email IN (?)", bun.In(emails)).Exec(ctx)
		if err != nil {
			logger.Error("failed to create test location users", zap.Error(err))
			return errors.Wrap(err, "failed to create test location users")
		}
		return nil
	}, nil
}
