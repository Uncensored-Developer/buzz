package repository

import (
	"context"
	"github.com/uptrace/bun"
	"time"
)

// SelectCriteria is a function type that takes a pointer to bun.SelectQuery as input and returns the modified query.
// It is used as a criteria function in repository methods to filter the entities based on specific conditions.
type SelectCriteria func(*bun.SelectQuery) *bun.SelectQuery

// IRepository is an interface that represents a generic repository for saving and retrieving entities.
// It provides methods for saving an entity, finding one entity based on criteria, and finding multiple entities based on criteria.
// The type parameter T represents the entity type that the repository operates on.
type IRepository[T any] interface {
	Save(context.Context, *T) error
	FindOne(context.Context, ...SelectCriteria) (T, error)
	FindAll(context.Context, ...SelectCriteria) ([]T, error)
	Delete(context.Context, *T) error
	Update(context.Context, *T) error
}

// ISimpleCacheManager is an interface for a simple cache repository that provides
// methods for setting a value with a key, retrieving a value using a key and deleting a value.
type ISimpleCacheManager interface {
	Set(ctx context.Context, key string, value string, ttl time.Duration) error
	Get(ctx context.Context, key string) (string, error)
	Delete(ctx context.Context, key string) error
}
