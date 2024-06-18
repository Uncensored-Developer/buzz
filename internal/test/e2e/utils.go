package e2e

import (
	"context"
	"github.com/Uncensored-Developer/buzz/internal/users/features"
	"github.com/Uncensored-Developer/buzz/internal/users/models"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"net/http"
	"time"
)

var ErrEndpointTimeout = errors.New("timeout request while waiting for endpoint/server.")

// WaitForReady waits for an endpoint/server to be ready by periodically sending GET requests until a successful
// response is received or the timeout expires.
//
// It takes a context, a Logger, a timeout duration, and an endpoint URL as input parameters.
// The context is used to provide cancellation support and to check for timeouts.
// The Logger is used to log any errors encountered during the process.
// The timeout is the maximum duration during which the function will wait for the endpoint/server to be ready.
// The endpoint URL is the URL of the endpoint/server that is being waited upon.
//
// The function returns an error if the endpoint/server is not ready within the specified timeout duration or
// if any errors occur during the waiting process.
// If the endpoint/server becomes ready within the timeout duration, the function returns nil.
func WaitForReady(
	ctx context.Context,
	logger *zap.Logger,
	timeout time.Duration,
	endpoint string,
) error {
	client := http.Client{}
	startTime := time.Now()

	for {
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
		if err != nil {
			return errors.Wrap(err, "create request failed")
		}

		res, err := client.Do(req)
		if err != nil {
			logger.Info("Error making request", zap.Error(err))
			continue
		}
		if res.StatusCode == http.StatusOK {
			logger.Info("Endpoint/Server is ready!")
			res.Body.Close()
			return nil
		}
		res.Body.Close()

		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			if time.Since(startTime) >= timeout {
				return ErrEndpointTimeout
			}
			// Wait a little before each checks
			time.Sleep(250 * time.Millisecond)
		}
	}
}

// CreateUser creates a new user by signing up with the specified information.
// It takes a context, date of birth, email, password, and gender as input parameters.
// The context is used to provide cancellation support.
// The date of birth is the user's date of birth, If empty a random date in the past would be generated.
// The email is the user's email address, if empty a random email address will be generated.
// The password is the user's password.
// The gender is the user's gender.
// The function returns a pointer to the created user and an error.
// If any error occurs during the user signup process, an error is returned along with a nil pointer.
// If the authentication service fails to initialize, an error is returned with a nil pointer.
func CreateUser(ctx context.Context, dob time.Time, email, password, gender string) (*models.User, error) {
	authService, err := features.InitializeAuthenticationService()
	if err != nil {
		return nil, errors.Wrap(err, "init auth service failed")
	}
	if email == "" {
		email = gofakeit.Email()
	}
	name := gofakeit.Name()
	if dob.IsZero() {
		dob = gofakeit.PastDate()
	}

	user, err := authService.SignUp(ctx, dob, name, email, password, gender)
	if err != nil {
		return nil, errors.Wrap(err, "user signup failed")
	}
	return &user, nil
}
