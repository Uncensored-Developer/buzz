package features_test

import (
	"context"
	"fmt"
	"github.com/Uncensored-Developer/buzz/internal/datastore"
	data2 "github.com/Uncensored-Developer/buzz/internal/matches/data"
	"github.com/Uncensored-Developer/buzz/internal/matches/features"
	"github.com/Uncensored-Developer/buzz/internal/test/e2e"
	"github.com/Uncensored-Developer/buzz/internal/users/data"
	features2 "github.com/Uncensored-Developer/buzz/internal/users/features"
	"github.com/Uncensored-Developer/buzz/internal/users/models"
	"github.com/Uncensored-Developer/buzz/pkg/cache"
	"github.com/Uncensored-Developer/buzz/pkg/config"
	"github.com/Uncensored-Developer/buzz/pkg/db"
	"github.com/Uncensored-Developer/buzz/pkg/logger"
	"github.com/Uncensored-Developer/buzz/pkg/migrate"
	"github.com/Uncensored-Developer/buzz/pkg/repository"
	"github.com/Uncensored-Developer/buzz/pkg/testcontainer"
	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"
	"os"
	"testing"
	"time"
)

type matchServiceTestSuite struct {
	suite.Suite
	ctx          context.Context
	logger       *zap.Logger
	matchService *features.MatchService
	matchesRepo  data2.IMatchesRepository
	userRepo     data.IUserRepository
	userOne      *models.User
	userTwo      *models.User
	cacheManager repository.ISimpleCacheManager

	testDatabase  *testcontainer.TestDatabase
	cacheDatabase *testcontainer.TestCacheDatabase
}

func (m *matchServiceTestSuite) SetupSuite() {
	m.ctx = context.Background()
	m.logger = logger.NewLogger()
	var err error
	// Start up test database container
	m.testDatabase, err = testcontainer.NewTestDatabase(m.ctx, m.logger)
	m.Require().NoError(err)

	// Run migrations on test database
	err = migrate.Up(m.testDatabase.DSN, "db/migrations")
	m.Require().NoError(err)

	m.cacheDatabase, err = testcontainer.NewCacheDatabase(m.ctx, m.logger)
	m.Require().NoError(err)

	// This is necessary for the config.LoadConfig to pick this from the environment variables
	err = os.Setenv("DATABASE_URL", m.testDatabase.DSN)
	m.Require().NoError(err)
	err = os.Setenv("REDIS_URL", m.cacheDatabase.DSN)
	m.Require().NoError(err)

	cfg, err := config.LoadConfig()
	m.Require().NoError(err)

	bunDb, err := db.Connect(m.testDatabase.DSN)
	m.Require().NoError(err)

	m.userRepo = data.NewUserRepository(bunDb)
	m.matchesRepo = data2.NewMatchesRepository(bunDb)
	m.cacheManager = cache.NewRedisManager(cfg)
	uow := datastore.NewUnitOfWorkDatastore(bunDb)

	m.matchService = features.NewMatchService(m.userRepo, m.cacheManager, uow, cfg, m.logger)

}

// TearDownSuite tears down the test suite by performing the following actions:
// 1. Unset the "DATABASE_URL" environment variable.
// 2. Unset the "REDIS_URL" environment variable.
// 3. Shutdown the testDatabase container.
// 4. Shutdown the cacheDatabase container.
//
// If any of the above actions fail, an error is logged using the logger.
func (m *matchServiceTestSuite) TearDownSuite() {

	err := os.Unsetenv("DATABASE_URL")
	if err != nil {
		m.logger.Error("unset DATABASE_URL env failed", zap.Error(err))
	}
	err = os.Unsetenv("REDIS_URL")
	if err != nil {
		m.logger.Error("unset REDIS_URL env failed", zap.Error(err))
	}

	err = m.testDatabase.Shutdown()
	if err != nil {
		m.logger.Error("failed to terminate database container", zap.Error(err))
	}

	err = m.cacheDatabase.Shutdown()
	if err != nil {
		m.logger.Error("failed to terminate cache database container", zap.Error(err))
	}
}

// SetupTest sets up the test data before each individual test case in the test suite.
// It creates two user records for testing purposes using the e2e.CreateUser function.
// The user records are assigned to the m.userOne and m.userTwo variables for later use in the test cases.
func (m *matchServiceTestSuite) SetupTest() {
	var err error
	m.userOne, err = e2e.CreateUser(m.ctx, time.Time{}, "", "FakeUserPassword", "M")
	m.Require().NoError(err)

	m.userTwo, err = e2e.CreateUser(m.ctx, time.Time{}, "", "FakeUserPassword", "M")
	m.Require().NoError(err)
}

// TearDownTest cleans up test data after each individual test case in the test suite.
// It deletes the user records used for testing purposes:
func (m *matchServiceTestSuite) TearDownTest() {
	err := m.userRepo.Delete(m.ctx, m.userOne)
	m.Require().NoError(err)

	err = m.userRepo.Delete(m.ctx, m.userTwo)
	m.Require().NoError(err)
}

func (m *matchServiceTestSuite) TestSwipe_InvalidSwiperIdReturnsError() {
	invalidSwiperID := int64(9999999)
	_, err := m.matchService.Swipe(m.ctx, invalidSwiperID, m.userTwo.ID, features.YesAction)
	m.Require().ErrorIs(err, features2.ErrUserNotFound)

	_, err = m.matchService.Swipe(m.ctx, m.userOne.ID, invalidSwiperID, features.YesAction)
	m.Require().ErrorIs(err, features2.ErrUserNotFound)
}

func (m *matchServiceTestSuite) TestSwipe_LikeButNoMatch() {
	userBefore, err := m.userRepo.FindOne(m.ctx, data.UserWithID(m.userTwo.ID))
	m.Require().NoError(err)

	match, err := m.matchService.Swipe(m.ctx, m.userOne.ID, m.userTwo.ID, features.YesAction)
	m.Assert().NoError(err)
	m.Assert().Empty(match)

	userAfter, err := m.userRepo.FindOne(m.ctx, data.UserWithID(m.userTwo.ID))
	m.Require().NoError(err)

	// Check if the liked user's like count was increased
	m.Assert().Equal(userBefore.LikesCount+1, userAfter.LikesCount)

	// Check if swipe action was actually saved to cache
	key := fmt.Sprintf("%d.%s.%d", m.userOne.ID, features.YesAction, m.userTwo.ID)
	_, err = m.cacheManager.Get(m.ctx, key)
	m.Assert().NoError(err)
}

func (m *matchServiceTestSuite) TestSwipe_LikeWithMatch() {
	userTwoBefore, err := m.userRepo.FindOne(m.ctx, data.UserWithID(m.userTwo.ID))
	m.Require().NoError(err)
	userOneBefore, err := m.userRepo.FindOne(m.ctx, data.UserWithID(m.userOne.ID))
	m.Require().NoError(err)

	// userTwo first LIKE swipe userOne
	_, err = m.matchService.Swipe(m.ctx, m.userTwo.ID, m.userOne.ID, features.YesAction)
	m.Assert().NoError(err)

	// The userOne LIKE swipes userTwo
	match, err := m.matchService.Swipe(m.ctx, m.userOne.ID, m.userTwo.ID, features.YesAction)
	m.Assert().NoError(err)
	m.Assert().NotEmpty(match)
	m.Assert().NotZero(match.ID)

	userTwoAfter, err := m.userRepo.FindOne(m.ctx, data.UserWithID(m.userTwo.ID))
	m.Require().NoError(err)
	userOneAfter, err := m.userRepo.FindOne(m.ctx, data.UserWithID(m.userOne.ID))
	m.Require().NoError(err)

	// Check if match was saved to matched repository
	gotMatch, err := m.matchesRepo.FindOne(m.ctx, data2.MatchWithID(match.ID))
	m.Assert().NoError(err)
	m.Assert().Equal(gotMatch, match)

	// Check if swipe action was actually deleted from the cache
	key := fmt.Sprintf("%d.%s.%d", m.userTwo.ID, features.YesAction, m.userOne.ID)
	err = m.cacheManager.Delete(m.ctx, key)
	m.Assert().NoError(err)

	// Check if the liked user's like count was increased
	m.Assert().Equal(userOneBefore.LikesCount+1, userOneAfter.LikesCount)
	m.Assert().Equal(userTwoBefore.LikesCount+1, userTwoAfter.LikesCount)
}

func TestMatchService(t *testing.T) {
	suite.Run(t, new(matchServiceTestSuite))
}
