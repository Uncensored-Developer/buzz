package utils_test

import (
	"github.com/Uncensored-Developer/buzz/pkg/utils"
	"github.com/pkg/errors"
	"math"
	"testing"
)

func TestConvertDatabaseUrlToDSN(t *testing.T) {
	testCases := map[string]struct {
		url           string
		expectedDSN   string
		expectErr     bool
		expectedError string
	}{
		"valid mysql URL": {
			url:         "mysql://user:password@localhost:3306/dbname",
			expectedDSN: "user:password@tcp(localhost:3306)/dbname",
			expectErr:   false,
		},
		"valid postgresql URL": {
			url:         "postgres://user:password@localhost:3306/dbname",
			expectedDSN: "user:password@tcp(localhost:3306)/dbname",
			expectErr:   false,
		},
		"missing host": {
			url:           "mysql://user:password@/dbname",
			expectedDSN:   "",
			expectErr:     true,
			expectedError: "invalid database URL: missing host",
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			dsn, err := utils.ConvertDatabaseUrlToDSN(tc.url)
			if tc.expectErr {
				if err == nil {
					t.Fatalf("expected error but got none")
				}

				if errors.Cause(err).Error() != tc.expectedError {
					t.Fatalf("expected error message '%s' but got '%v'", tc.expectedError, errors.Cause(err))
				}
			} else {
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}

				if dsn != tc.expectedDSN {
					t.Fatalf("expected DSN '%s' but got '%s'", tc.expectedDSN, dsn)
				}
			}
		})
	}
}

func TestDistanceBetween(t *testing.T) {
	testCases := map[string]struct {
		lat1, long1  float64
		lat2, long2  float64
		expectedDist float64
	}{
		"San Francisco to Los Angeles": {
			lat1:         37.7749,
			long1:        -122.4194,
			lat2:         34.0522,
			long2:        -118.2437,
			expectedDist: 559.0, // Approximate distance in KM
		},
		"New York to Los Angeles": {
			lat1:         40.7128,
			long1:        -74.0060,
			lat2:         34.0522,
			long2:        -118.2437,
			expectedDist: 3935.0,
		},
		"London to Paris": {
			lat1:         51.5074,
			long1:        -0.1278,
			lat2:         48.8566,
			long2:        2.3522,
			expectedDist: 344.0,
		},
		"Tokyo to Sydney": {
			lat1:         35.6895,
			long1:        139.6917,
			lat2:         -33.8688,
			long2:        151.2093,
			expectedDist: 7826.0,
		},
		"UserA to UserB": {
			lat1:         51.2725887,
			long1:        0.5026768,
			lat2:         50.96284649,
			long2:        -0.12981616,
			expectedDist: 55.0,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			got := utils.DistanceBetween(tc.lat1, tc.long1, tc.lat2, tc.long2)
			if math.Abs(got-tc.expectedDist) > 1.0 { // Allowing 1 km error
				t.Errorf(
					"DistanceBetween(%f, %f, %f, %f) = %f; want %f",
					tc.lat1, tc.long1, tc.lat2, tc.long2, got, tc.expectedDist)
			}
		})
	}
}
