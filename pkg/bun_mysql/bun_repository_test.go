package bun_mysql_test

import (
	"context"
	"github.com/Uncensored-Developer/buzz/pkg/bun_mysql"
	"github.com/Uncensored-Developer/buzz/pkg/db"
	"github.com/Uncensored-Developer/buzz/pkg/testcontainer"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/uptrace/bun"
	"testing"
)

type testUser struct {
	ID       int
	Email    string
	Password string
}

func TestBunRepository(t *testing.T) {
	ctx := context.Background()
	testDb, err := testcontainer.InitializeTestDatabase(ctx)
	require.NoError(t, err)
	defer testDb.Shutdown()

	bunDB, err := db.Connect(testDb.DSN)
	require.NoError(t, err)

	_, err = bunDB.NewCreateTable().Model(&testUser{}).Exec(ctx)
	require.NoError(t, err)

	testUserBunRepository := bun_mysql.NewBunRepository[testUser](bunDB)
	testUsers := []testUser{
		{Email: "user1@buzz.com", Password: "password1"},
		{Email: "user2@buzz.com", Password: "password2"},
	}

	t.Run("save model", func(t *testing.T) {
		users, err := testUserBunRepository.FindAll(ctx)
		assert.NoError(t, err)
		assert.Equal(t, 0, len(users))

		for i := range testUsers {
			err = testUserBunRepository.Save(ctx, &testUsers[i])
			assert.NoError(t, err)
		}
	})

	t.Run("find all", func(t *testing.T) {
		users, err := testUserBunRepository.FindAll(ctx)
		assert.NoError(t, err)
		assert.ElementsMatch(t, testUsers, users)
	})

	t.Run("find all with filters", func(t *testing.T) {
		c := func(q *bun.SelectQuery) *bun.SelectQuery {
			return q.Where("email = ?", testUsers[0].Email)
		}

		users, err := testUserBunRepository.FindAll(ctx, c)
		assert.NoError(t, err)
		assert.Equal(t, 1, len(users))
		assert.Equal(t, testUsers[0], users[0])
	})

	t.Run("find one", func(t *testing.T) {
		user, err := testUserBunRepository.FindOne(ctx)
		assert.NoError(t, err)
		assert.Equal(t, testUsers[0], user)
	})

	t.Run("find one with filters", func(t *testing.T) {
		c := func(q *bun.SelectQuery) *bun.SelectQuery {
			return q.Where("email = ?", testUsers[1].Email)
		}

		user, err := testUserBunRepository.FindOne(ctx, c)
		assert.NoError(t, err)
		assert.Equal(t, testUsers[1], user)
	})

}
