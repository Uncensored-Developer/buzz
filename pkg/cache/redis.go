package cache

import (
	"context"
	"fmt"
	"github.com/Uncensored-Developer/buzz/pkg/config"
	"github.com/pkg/errors"
	"github.com/redis/go-redis/v9"
	"sync"
	"time"
)

const KeyPrefix = "BUZZ_APP_"

type RedisManager struct {
	client *redis.Client
}

var singleton *RedisManager
var once sync.Once

func NewRedisManager(cfg *config.Config) *RedisManager {
	once.Do(func() {
		opts, err := redis.ParseURL(cfg.RedisURL)
		if err != nil {
			panic(fmt.Sprintf("Redis Opt failed: %v", err))
		}
		client := redis.NewClient(opts)
		singleton = &RedisManager{
			client: client,
		}
	})
	return singleton
}

func (r *RedisManager) Set(ctx context.Context, key string, value string, ttl time.Duration) error {
	err := r.client.Set(ctx, KeyPrefix+key, value, ttl).Err()
	if err != nil {
		return errors.Wrap(err, "Redis SET failed")
	}
	return nil
}

func (r *RedisManager) Get(ctx context.Context, key string) (string, error) {
	val, err := r.client.Get(ctx, KeyPrefix+key).Result()
	if err != nil {
		return "", errors.Wrap(err, "Redis GET failed")
	}
	return val, nil
}

func (r *RedisManager) Delete(ctx context.Context, key string) error {
	_, err := r.client.Del(ctx, KeyPrefix+key).Result()
	if err != nil {
		return errors.Wrap(err, "Redis DEL failed")
	}
	return nil
}
