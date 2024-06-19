package utils

import (
	"fmt"
	"github.com/pkg/errors"
	"math"
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

// DistanceBetween calculates the distance between two points on the Earth's surface using the Haversine formula.
// The distance is returned in kilometers.
func DistanceBetween(lat1, lon1, lat2, lon2 float64) float64 {
	const R = 6371 // Earth radius in kilometers

	// Convert latitude and longitude from degrees to radians
	lat1Rad := lat1 * math.Pi / 180
	lon1Rad := lon1 * math.Pi / 180
	lat2Rad := lat2 * math.Pi / 180
	lon2Rad := lon2 * math.Pi / 180

	// Differences in coordinates
	dLat := lat2Rad - lat1Rad
	dLon := lon2Rad - lon1Rad

	// Haversine formula
	a := math.Sin(dLat/2)*math.Sin(dLat/2) +
		math.Cos(lat1Rad)*math.Cos(lat2Rad)*
			math.Sin(dLon/2)*math.Sin(dLon/2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	return R * c // Distance in kilometers
}
