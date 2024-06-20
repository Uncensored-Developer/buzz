package cache_test

import (
	"context"
	"github.com/Uncensored-Developer/buzz/pkg/cache"
	"github.com/Uncensored-Developer/buzz/pkg/config"
	"github.com/Uncensored-Developer/buzz/pkg/logger"
	"github.com/Uncensored-Developer/buzz/pkg/testcontainer"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"
	"os"
	"testing"
	"time"
)

type redisManagerTestSuite struct {
	suite.Suite
	ctx           context.Context
	manager       *cache.RedisManager
	redisURL      string
	redisDatabase *testcontainer.TestCacheDatabase
	redisClient   *redis.Client
	logger        *zap.Logger
}

func (r *redisManagerTestSuite) SetupSuite() {
	ctx := context.Background()
	zapLogger := logger.NewLogger()
	redisContainer, err := testcontainer.NewCacheDatabase(ctx, zapLogger)
	r.Require().NoError(err)
	r.redisDatabase = redisContainer
	r.logger = zapLogger

	err = os.Setenv("BUZZ_REDIS_URL", redisContainer.DSN)
	r.Require().NoError(err)

	r.redisURL = redisContainer.DSN

	cfg, err := config.LoadConfig()
	r.Require().NoError(err)
	redisManager := cache.NewRedisManager(cfg)
	r.manager = redisManager

	r.redisClient, err = getRedisClient(r.redisURL)
	r.Require().NoError(err)

	r.ctx = ctx
}

func (r *redisManagerTestSuite) TearDownSuite() {
	err := os.Unsetenv("BUZZ_REDIS_URL")
	if err != nil {
		r.logger.Error("unset BUZZ_REDIS_URL env failed", zap.Error(err))
	}
	err = r.redisDatabase.Shutdown()
	if err != nil {
		r.logger.Error("failed to terminate redis database container", zap.Error(err))
	}
}

// TestSave tests the Save method of the RedisManager.
// It tests if the method save a value to the database with the right prefix to the key
func (r *redisManagerTestSuite) TestSave() {
	key := "testKey"
	value := "test-value"

	err := r.manager.Set(r.ctx, key, value, time.Minute)
	r.Require().NoError(err)

	gotValue, err := r.redisClient.Get(r.ctx, cache.KeyPrefix+key).Result()
	r.Require().NoError(err)

	r.Assert().Equal(gotValue, value)
}

// TestGet tests the Get method of the RedisManager.
// It tests if the method retrieves the correct value from the database using the provided key accounting for the key prefix.
func (r *redisManagerTestSuite) TestGet() {
	key := cache.KeyPrefix + "testKey2"
	value := "test_value"

	err := r.redisClient.Set(r.ctx, key, value, time.Minute).Err()
	r.Require().NoError(err)

	gotValue, err := r.manager.Get(r.ctx, "testKey2")
	r.Require().NoError(err)

	r.Assert().Equal(gotValue, value)
}

// TestGetNonExistingKey tests the Get method of the RedisManager.
// It tests if the method returns an error when trying to get a non-existing key from the database.
func (r *redisManagerTestSuite) TestGetNonExistingKey() {
	_, err := r.manager.Get(r.ctx, "nonExistingKey321")
	r.Require().Error(err)
}

func (r *redisManagerTestSuite) TestDelete() {
	key := cache.KeyPrefix + "testKey3"
	value := "test_value"

	err := r.redisClient.Set(r.ctx, key, value, time.Minute).Err()
	r.Require().NoError(err)

	err = r.manager.Delete(r.ctx, "testKey3")
	r.Require().NoError(err)

	_, err = r.redisClient.Get(r.ctx, key).Result()
	r.Require().Error(err)
}

func getRedisClient(url string) (*redis.Client, error) {
	opts, err := redis.ParseURL(url)
	if err != nil {
		return nil, err
	}
	return redis.NewClient(opts), nil
}

func TestRedisManager(t *testing.T) {
	suite.Run(t, new(redisManagerTestSuite))
}
