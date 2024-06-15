package datastore

import (
	"context"
	"database/sql"
	"github.com/Uncensored-Developer/buzz/internal/users/data"
	"github.com/uptrace/bun"
)

// IUnitOfWorkDatastore is an interface that represents a unit of work for a data store.
// It provides access to all data stores in a Unit-Of-Work.
// All data changes done will be executed in a database transaction
type IUnitOfWorkDatastore interface {
	UsersRepository() data.IUserRepository
}

type UnitOfWorkBlock func(store IUnitOfWorkDatastore) error

type IUnitOfWork interface {
	Do(ctx context.Context, store UnitOfWorkBlock) error
}

type unitOfWorkDataStore struct {
	usersRepository data.IUserRepository
}

func (u *unitOfWorkDataStore) UsersRepository() data.IUserRepository {
	return u.usersRepository
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
			usersRepository: data.NewUserRepository(tx),
		}
		return fn(s)
	})
}
