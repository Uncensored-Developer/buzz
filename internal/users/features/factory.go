package features

import (
	"github.com/Uncensored-Developer/buzz/internal/users/data"
	"github.com/Uncensored-Developer/buzz/pkg/authentication"
	"github.com/Uncensored-Developer/buzz/pkg/config"
	"github.com/Uncensored-Developer/buzz/pkg/db"
	"github.com/Uncensored-Developer/buzz/pkg/hash"
	"github.com/Uncensored-Developer/buzz/pkg/logger"
	"github.com/pkg/errors"
)

func InitializeAuthenticationService() (*AuthenticationService, error) {
	cfg, err := config.LoadConfig()
	if err != nil {
		return nil, errors.Wrap(err, "load config failed")
	}
	zapLogger := logger.NewLogger()

	passwordHasher := hash.NewSHA1Hasher(cfg.PasswordHasherSalt)
	manager, err := authentication.NewManager(cfg.JwtKey)
	if err != nil {
		return nil, errors.Wrap(err, "token manager error")
	}

	bunDB, err := db.Connect(cfg.DatabaseURL)
	if err != nil {
		return nil, errors.Wrap(err, "database connection failed")
	}
	userRepo := data.NewUserRepository(bunDB)
	return NewAuthenticationService(
		passwordHasher,
		manager,
		userRepo,
		cfg,
		zapLogger,
	), nil
}
