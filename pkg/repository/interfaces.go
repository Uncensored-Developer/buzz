package repository

import (
	"context"
	"github.com/uptrace/bun"
)

type SelectCriteria func(*bun.SelectQuery) *bun.SelectQuery

type IRepository[T any] interface {
	Save(context.Context, *T) error
	FindOne(context.Context, ...SelectCriteria) (T, error)
	FindAll(context.Context, ...SelectCriteria) ([]T, error)
}

type ISimpleCacheRepository interface {
	Set(ctx context.Context, key string, value string) error
	Get(ctx context.Context, key string) (string, error)
}
