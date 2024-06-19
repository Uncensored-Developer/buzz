package bun_mysql

import (
	"context"
	"database/sql"
	"github.com/Uncensored-Developer/buzz/pkg/repository"
	"github.com/pkg/errors"
	"github.com/uptrace/bun"
)

var ErrRowNotFound = errors.New("row not found")

// BunRepository is a generic repository implementation that uses bun as the underlying database ORM.
type BunRepository[T any] struct {
	DB bun.IDB
}

func NewBunRepository[T any](db bun.IDB) repository.IRepository[T] {
	return &BunRepository[T]{DB: db}
}

// Save saves the provided model in the database.
// It accepts a context and a pointer to a model of type T.
// It returns an error if there was an issue with the save operation.
func (b *BunRepository[T]) Save(ctx context.Context, model *T) error {
	_, err := b.DB.NewInsert().Model(model).Returning("*").Exec(ctx)
	return err
}

// FindOne finds a single entity based on the provided selection criteria.
// It accepts a context and a variadic list of 'repository.SelectCriteria' functions that modify the SelectQuery.
// It returns the found entity of type T and an error, if any.
// If no entity is found, it returns the zero-value of T and an error containing the message "could not find entity."
func (b *BunRepository[T]) FindOne(ctx context.Context, filters ...repository.SelectCriteria) (T, error) {
	var row T

	q := b.DB.NewSelect().Model(&row)
	for i := range filters {
		q.Apply(filters[i])
	}

	err := q.Limit(1).Scan(ctx)
	if errors.Is(err, sql.ErrNoRows) {
		return row, ErrRowNotFound
	}
	return row, err
}

// FindAll retrieves all models from the database.
// It accepts a context and optional filters to be applied when querying the database.
// The filters should be of type 'repository.SelectCriteria'.
// It returns a slice of models of type T and an error if there was an issue with the query.
func (b *BunRepository[T]) FindAll(ctx context.Context, filters ...repository.SelectCriteria) ([]T, error) {
	var rows []T

	q := b.DB.NewSelect().Model(&rows)
	for i := range filters {
		q.Apply(filters[i])
	}

	err := q.Scan(ctx)
	return rows, err
}

func (b *BunRepository[T]) Delete(ctx context.Context, model *T) error {
	_, err := b.DB.NewDelete().Model(model).WherePK().Exec(ctx)
	return err
}

func (b *BunRepository[T]) Update(ctx context.Context, model *T) error {
	_, err := b.DB.NewUpdate().Model(model).WherePK().Returning("*").Exec(ctx)
	return err
}
