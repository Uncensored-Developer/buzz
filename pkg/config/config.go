package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/pkg/errors"
)

type Config struct {
	Host               string `env:"BUZZ_HTTP_HOST" env-default:"localhost"`
	Port               string `env:"BUZZ_HTTP_PORT" env-default:"8010"`
	Debug              bool   `env:"BUZZ_DEBUG" env-default:"true"`
	DatabaseURL        string `env:"BUZZ_DATABASE_URL"`
	RedisURL           string `env:"BUZZ_REDIS_URL"`
	JwtKey             string `env:"BUZZ_JWT_KEY" env-default:"fakeJwtkey"`
	PasswordHasherSalt string `env:"BUZZ_PASSWORD_HASHER_SALT" env-default:"fakeHasherSalt"`
	FakeUserPassword   string `env:"BUZZ_FAKE_USER_PASSWORD" env-default:"password123"`
	SeedUsers          bool   `env:"BUZZ_SEED_USERS" env-default:"false"`

	// For good balance between precision (for a dating app use case) and performance
	H3Resolution int `env:"H3_RESOLUTION" env-default:"9"`
}

func LoadConfig() (*Config, error) {
	cfg := &Config{}
	err := cleanenv.ReadEnv(cfg)
	if err != nil {
		return nil, errors.Wrap(err, "could not read config")
	}
	return cfg, nil
}
