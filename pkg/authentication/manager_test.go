package authentication_test

import (
	"github.com/Uncensored-Developer/buzz/pkg/authentication"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"strconv"
	"testing"
	"time"
)

// TestManager is a test function that tests the functionality of the TestManager.
// It simply just tests that NewToken() and Parse() works as expected
// i.e. token gotten from the userId can be parsed back to the same userId.
func TestManager(t *testing.T) {
	manager, err := authentication.NewManager("testSigningKey")
	require.NoError(t, err)

	var userId int64 = 100

	token, err := manager.NewToken(userId, time.Duration(24*time.Hour))
	require.NoError(t, err)

	assert.NotEqual(t, strconv.FormatInt(userId, 10), token)

	got, err := manager.Parse(token)
	require.NoError(t, err)
	assert.Equal(t, userId, got)
}

// TODO: Tests for errors
