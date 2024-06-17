package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/pkg/errors"
)

type Config struct {
	Host               string `env:"HTTP_HOST"`
	Port               string `env:"HTTP_PORT"`
	Debug              bool   `env:"DEBUG" env-default:"true"`
	DatabaseURL        string `env:"DATABASE_URL"`
	JwtKey             string `env:"JWT_KEY" env-default:"fakeJwtkey"`
	PasswordHasherSalt string `env:"PASSWORD_HASHER_SALT" env-default:"fakeHasherSalt"`
	FakeUserPassword   string `env:"FAKE_USER_PASSWORD" env-default:"password123"`
}

func LoadConfig() (*Config, error) {
	cfg := &Config{}
	err := cleanenv.ReadEnv(cfg)
	if err != nil {
		return nil, errors.Wrap(err, "could not read config")
	}
	return cfg, nil
}
