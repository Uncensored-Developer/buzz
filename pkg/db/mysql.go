package db

import (
	"database/sql"
	"github.com/Uncensored-Developer/buzz/pkg/utils"
	_ "github.com/go-sql-driver/mysql"
	"github.com/pkg/errors"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/mysqldialect"
)

func Connect(dbURL string) (*bun.DB, error) {
	dsn, err := utils.ConvertDatabaseUrlToDSN(dbURL)
	if err != nil {
		return nil, errors.Wrap(err, "url conversion failed")
	}
	sqlDb, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, errors.Wrap(err, "could not open database connection")
	}
	db := bun.NewDB(sqlDb, mysqldialect.New())
	return db, nil
}
