package data

import (
	"context"
	"github.com/Uncensored-Developer/buzz/internal/users/models"
	"github.com/Uncensored-Developer/buzz/pkg/db"
	"github.com/Uncensored-Developer/buzz/pkg/migrate"
	"github.com/Uncensored-Developer/buzz/pkg/testcontainer"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"sync"
	"testing"
	"time"
)

// TestUserRepository_IncrementLikes is a test function to increment the LikesCount of a user in the UserRepository.
// It initializes a test database, migrates the necessary schema, and creates a user with an initial LikesCount of 1.
// It then increments the LikesCount multiple times using goroutines to simulate concurrent requests.
// Finally, it checks if the LikesCount has been incremented correctly.
func TestUserRepository_IncrementLikes(t *testing.T) {
	ctx := context.Background()
	testDb, err := testcontainer.InitializeTestDatabase(ctx)
	require.NoError(t, err)
	defer testDb.Shutdown()

	err = migrate.Up(testDb.DSN, "db/migrations")
	require.NoError(t, err)

	bunDb, err := db.Connect(testDb.DSN)
	require.NoError(t, err)

	userRepo := NewUserRepository(bunDb)

	email := "test.user@buzz.com"
	user := models.User{
		Name:       "John Doe",
		Email:      email,
		Password:   "hashedPassword",
		Gender:     "M",
		Dob:        time.Now(),
		LikesCount: 1,
	}
	err = userRepo.Save(ctx, &user)
	require.NoError(t, err)

	loops := 5

	gotUserBefore, err := userRepo.FindOne(ctx, UserWithEmail(email))
	require.NoError(t, err)
	countBefore := gotUserBefore.LikesCount
	var wg sync.WaitGroup
	wg.Add(loops)
	for i := 0; i < loops; i++ {
		go func(wg *sync.WaitGroup) {
			defer wg.Done()
			err := userRepo.IncrementLikes(ctx, gotUserBefore.ID, 1)
			require.NoError(t, err)
		}(&wg)
	}
	wg.Wait()

	gotUserAfter, err := userRepo.FindOne(ctx, UserWithEmail(email))
	require.NoError(t, err)
	countAfter := gotUserAfter.LikesCount
	assert.Equal(t, countBefore+loops, countAfter)
}
