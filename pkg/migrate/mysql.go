package migrate

import (
	"database/sql"
	"fmt"
	"github.com/Uncensored-Developer/buzz/pkg/testcontainer"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/mysql"
	"github.com/pkg/errors"
	"path/filepath"
	"runtime"
)

// Up applies all available migrations to the specified database.
// It opens a connection to the database using the provided DSN (Data Source Name),
// loads the migrations from the given source path, and applies them to the database.
// The source path is determined by the project's root directory relative to the caller's file location.
// The database driver is set as MySQL.
// If any errors occur during the process, they will be returned.
func Up(dsn, migrationsSource string) error {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return errors.Wrap(err, "could not open database connection")
	}

	driver, err := mysql.WithInstance(db, &mysql.Config{})
	if err != nil {
		return err
	}

	// find out the absolute path to this file
	// it'll be used to determine the project's root path
	_, callerPath, _, _ := runtime.Caller(0) // nolint:dogsled

	// look for migrations source starting from project's root dir
	sourceURL := fmt.Sprintf(
		"file://%s/../../%s",
		filepath.ToSlash(filepath.Dir(callerPath)),
		filepath.ToSlash(migrationsSource),
	)

	migration, err := migrate.NewWithDatabaseInstance(sourceURL, testcontainer.DbName, driver)
	if err != nil {
		return err
	}
	return migration.Up()
}
