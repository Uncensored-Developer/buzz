package e2e

import (
	"context"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"net/http"
	"time"
)

var ErrEndpointTimeout = errors.New("timeout request while waiting for endpoint/server.")

// WaitForReady waits for an endpoint/server to be ready by periodically sending GET requests until a successful
// response is received or the timeout expires.
//
// It takes a context, a logger, a timeout duration, and an endpoint URL as input parameters.
// The context is used to provide cancellation support and to check for timeouts.
// The logger is used to log any errors encountered during the process.
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
