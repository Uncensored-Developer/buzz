package datastore

import (
	"context"
	"database/sql"
	data2 "github.com/Uncensored-Developer/buzz/internal/matches/data"
	"github.com/Uncensored-Developer/buzz/internal/users/data"
	"github.com/uptrace/bun"
)

// IUnitOfWorkDatastore is an interface that represents a unit of work for a data store.
// It provides access to all data stores in a Unit-Of-Work.
// All data changes done will be executed in a database transaction
type IUnitOfWorkDatastore interface {
	UsersRepository() data.IUserRepository
	MatchesRepository() data2.IMatchesRepository
}

type UnitOfWorkBlock func(store IUnitOfWorkDatastore) error

type IUnitOfWork interface {
	Do(ctx context.Context, store UnitOfWorkBlock) error
}

type unitOfWorkDataStore struct {
	usersRepository   data.IUserRepository
	matchesRepository data2.IMatchesRepository
}

func (u *unitOfWorkDataStore) UsersRepository() data.IUserRepository {
	return u.usersRepository
}

func (u *unitOfWorkDataStore) MatchesRepository() data2.IMatchesRepository {
	return u.matchesRepository
}

type unitOfWork struct {
	conn *bun.DB
}

func NewUnitOfWorkDatastore(db *bun.DB) IUnitOfWork {
	return &unitOfWork{conn: db}
}

func (u *unitOfWork) Do(ctx context.Context, fn UnitOfWorkBlock) error {
	return u.conn.RunInTx(ctx, &sql.TxOptions{}, func(ctx context.Context, tx bun.Tx) error {
		s := &unitOfWorkDataStore{
			usersRepository:   data.NewUserRepository(tx),
			matchesRepository: data2.NewMatchesRepository(tx),
		}
		return fn(s)
	})
}
