package datastore_test

import (
	"context"
	"errors"
	"github.com/Uncensored-Developer/buzz/internal/datastore"
	models2 "github.com/Uncensored-Developer/buzz/internal/matches/models"
	"github.com/Uncensored-Developer/buzz/internal/users/data"
	"github.com/Uncensored-Developer/buzz/internal/users/models"
	"github.com/Uncensored-Developer/buzz/pkg/bun_mysql"
	"github.com/Uncensored-Developer/buzz/pkg/db"
	"github.com/Uncensored-Developer/buzz/pkg/migrate"
	"github.com/Uncensored-Developer/buzz/pkg/testcontainer"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestUnitOfWork(t *testing.T) {
	ctx := context.Background()

	testDb, err := testcontainer.InitializeTestDatabase(ctx)
	require.NoError(t, err)
	defer testDb.Shutdown()

	// migrations up
	err = migrate.Up(testDb.DSN, "db/migrations")
	require.NoError(t, err)

	bunDb, err := db.Connect(testDb.DSN)
	require.NoError(t, err)
	uow := datastore.NewUnitOfWorkDatastore(bunDb)
	usersRepo := data.NewUserRepository(bunDb)

	user := &models.User{
		Email:    "user1@buzz.com",
		Password: "password",
		Name:     "John Doe",
		Gender:   "M",
		Dob:      time.Date(1990, time.January, 1, 0, 0, 0, 0, time.UTC),
	}
	user2 := &models.User{
		Email:    gofakeit.Email(),
		Password: "password",
		Name:     "John Doe",
		Gender:   "M",
		Dob:      time.Date(1990, time.January, 1, 0, 0, 0, 0, time.UTC),
	}
	user3 := &models.User{
		Email:    gofakeit.Email(),
		Password: "password",
		Name:     "John Doe",
		Gender:   "M",
		Dob:      time.Date(1990, time.January, 1, 0, 0, 0, 0, time.UTC),
	}
	err = usersRepo.Save(ctx, user2)
	require.NoError(t, err)
	err = usersRepo.Save(ctx, user3)
	require.NoError(t, err)

	match := &models2.Match{
		UserOneID: user2.ID,
		UserTwoID: user2.ID,
	}

	t.Run("rollback on error", func(t *testing.T) {
		_, err := usersRepo.FindOne(ctx, data.UserWithEmail(user.Email))
		require.ErrorIs(t, err, bun_mysql.ErrRowNotFound)

		err = uow.Do(ctx, func(store datastore.IUnitOfWorkDatastore) error {
			err := store.UsersRepository().Save(ctx, user)
			require.NoError(t, err)

			err = store.MatchesRepository().Save(ctx, match)
			require.NoError(t, err)

			return errors.New("rollback error")
		})

		if assert.EqualError(t, err, "rollback error") {
			_, err := usersRepo.FindOne(ctx, data.UserWithEmail(user.Email))
			require.ErrorIs(t, err, bun_mysql.ErrRowNotFound)
		}
	})

	t.Run("rollback on panic with error", func(t *testing.T) {
		_, err := usersRepo.FindOne(ctx, data.UserWithEmail(user.Email))
		require.ErrorIs(t, err, bun_mysql.ErrRowNotFound)

		defer func() {
			p := recover()
			if assert.NotNil(t, p) {
				_, err := usersRepo.FindOne(ctx, data.UserWithEmail(user.Email))
				require.ErrorIs(t, err, bun_mysql.ErrRowNotFound)
			}
		}()

		err = uow.Do(ctx, func(store datastore.IUnitOfWorkDatastore) error {
			err := store.UsersRepository().Save(ctx, user)
			require.NoError(t, err)

			err = store.MatchesRepository().Save(ctx, match)
			require.NoError(t, err)

			panic(errors.New("rollback error"))
		})
	})

	t.Run("rollback on cancelled context", func(t *testing.T) {
		cancelCtx, cancel := context.WithCancel(ctx)

		_, err := usersRepo.FindOne(cancelCtx, data.UserWithEmail(user.Email))
		require.ErrorIs(t, err, bun_mysql.ErrRowNotFound)

		err = uow.Do(ctx, func(store datastore.IUnitOfWorkDatastore) error {
			err := store.UsersRepository().Save(cancelCtx, user)
			require.NoError(t, err)

			cancel()
			// This should actually return an error due to the canceled context
			return store.MatchesRepository().Save(cancelCtx, match)
		})

		if assert.EqualError(t, err, "context canceled") {
			_, err := usersRepo.FindOne(ctx, data.UserWithEmail(user.Email))
			require.ErrorIs(t, err, bun_mysql.ErrRowNotFound)
		}
	})

	t.Run("rollback on panic without error", func(t *testing.T) {
		_, err := usersRepo.FindOne(ctx, data.UserWithEmail(user.Email))
		require.ErrorIs(t, err, bun_mysql.ErrRowNotFound)

		defer func() {
			p := recover()
			if assert.NotNil(t, p) && assert.Equal(t, "rollback error", p) {
				_, err := usersRepo.FindOne(ctx, data.UserWithEmail(user.Email))
				require.ErrorIs(t, err, bun_mysql.ErrRowNotFound)
			}
		}()

		err = uow.Do(ctx, func(store datastore.IUnitOfWorkDatastore) error {
			err := store.UsersRepository().Save(ctx, user)
			require.NoError(t, err)

			err = store.MatchesRepository().Save(ctx, match)
			require.NoError(t, err)

			panic("rollback error")
		})
	})

	t.Run("commit on success", func(t *testing.T) {
		_, err := usersRepo.FindOne(ctx, data.UserWithEmail(user.Email))
		require.ErrorIs(t, err, bun_mysql.ErrRowNotFound)

		err = uow.Do(ctx, func(store datastore.IUnitOfWorkDatastore) error {
			err := store.UsersRepository().Save(ctx, user)
			require.NoError(t, err)

			err = store.MatchesRepository().Save(ctx, match)
			require.NoError(t, err)
			return nil
		})

		if assert.NoError(t, err) {
			gotUser, err := usersRepo.FindOne(ctx, data.UserWithEmail(user.Email))
			if assert.NoError(t, err) {
				user.ID = gotUser.ID
				user.CreatedAt = gotUser.CreatedAt
				user.UpdatedAt = gotUser.UpdatedAt
				assert.Equal(t, *user, gotUser)
			}
		}
	})
}
