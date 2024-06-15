package testcontainer

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/mysql"
	"go.uber.org/zap"
	"strings"
)

const DbUsername = "user"
const DbPassword = "password"
const DbName = "test_buzz"

type TestDatabase struct {
	container testcontainers.Container
	ctx       context.Context
	DSN       string
	logger    *zap.Logger
}

// NewTestDatabase creates a new instance of TestDatabase.
// It starts a MySQL test database container using testcontainers.RunContainer function and the given context and logger.
// The DSN (Data Source Name) for the MySQL test database is constructed using the provided DbUsername, DbPassword, mappedPort, and DbName.
// It returns a pointer to TestDatabase and an error if any.
//
// Example usage:
//
//	ctx := context.Background()
//	logger, err := zap.NewProduction()
//	if err != nil {
//	    log.Fatal(err)
//	}
//
//	testDB, err := NewTestDatabase(ctx, logger)
//	if err != nil {
//	    log.Fatal(err)
//	}
//	defer testDB.Shutdown()
//	// Use testDB instance for testing
func NewTestDatabase(ctx context.Context, logger *zap.Logger) (*TestDatabase, error) {
	dsnTemplate := "mysql://%s:%s@localhost:%s/%s"

	c, err := mysql.RunContainer(ctx,
		testcontainers.WithImage("mysql:8.4"),
		mysql.WithDatabase(DbName),
		mysql.WithUsername(DbUsername),
		mysql.WithPassword(DbPassword),
	)
	if err != nil {
		return nil, errors.Wrap(err, "failed to start mysql test db")
	}

	mappedPort, err := c.MappedPort(ctx, "3306/tcp")
	if err != nil {
		return nil, errors.Wrap(err, "failed to get port")
	}
	mappedPortStr := strings.Replace(string(mappedPort), "/tcp", "", 1)
	dsn := fmt.Sprintf(dsnTemplate, DbUsername, DbPassword, mappedPortStr, DbName)
	logger.Info("MYSQL test database container started successfully.")
	return &TestDatabase{
		container: c,
		ctx:       ctx,
		logger:    logger,
		DSN:       dsn,
	}, nil
}

func (t *TestDatabase) Shutdown() error {
	return t.container.Terminate(t.ctx)
}
