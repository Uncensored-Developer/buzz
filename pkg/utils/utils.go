package utils

import (
	"fmt"
	"github.com/pkg/errors"
	"net/url"
	"strings"
)

func ConvertDatabaseUrlToDSN(dbURL string) (string, error) {
	u, err := url.Parse(dbURL)
	if err != nil {
		return "", errors.Wrap(err, "URL parse failed")
	}

	user := u.User.Username()
	password, _ := u.User.Password()
	host := u.Host
	if host == "" {
		return "", errors.New("invalid database URL: missing host")
	}

	dbName := strings.TrimPrefix(u.Path, "/")
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s", user, password, host, dbName)
	return dsn, nil
}
