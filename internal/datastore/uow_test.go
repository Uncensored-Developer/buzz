package datastore_test

import (
	"context"
	"github.com/Uncensored-Developer/buzz/internal/datastore"
	"github.com/Uncensored-Developer/buzz/internal/users/data"
	"github.com/Uncensored-Developer/buzz/internal/users/models"
	"github.com/Uncensored-Developer/buzz/pkg/db"
	"github.com/Uncensored-Developer/buzz/pkg/migrate"
	"github.com/Uncensored-Developer/buzz/test/wire"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/uptrace/bun"
	"testing"
	"time"
)

func TestUnitOfWork(t *testing.T) {
	ctx := context.Background()

	testDb, err := wire.InitializeTestDatabase(ctx)
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

	t.Run("commit on success", func(t *testing.T) {
		err := uow.Do(ctx, func(store datastore.IUnitOfWorkDatastore) error {
			err := store.UsersRepository().Save(ctx, user)
			require.NoError(t, err)
			return nil
		})

		if assert.NoError(t, err) {
			c := func(q *bun.SelectQuery) *bun.SelectQuery {
				return q.Where("email = ?", user.Email)
			}
			gotUser, err := usersRepo.FindOne(ctx, c)
			if assert.NoError(t, err) {
				user.ID = gotUser.ID
				user.CreatedAt = gotUser.CreatedAt
				user.UpdatedAt = gotUser.UpdatedAt
				assert.Equal(t, *user, gotUser)
			}
		}
	})
}
