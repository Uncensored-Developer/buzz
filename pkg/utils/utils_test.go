package utils_test

import (
	"github.com/Uncensored-Developer/buzz/pkg/utils"
	"github.com/pkg/errors"
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
