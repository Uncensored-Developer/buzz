package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/pkg/errors"
)

type Config struct {
	Host  string `yaml:"host" env:"HTTP_HOST"`
	Port  string `yaml:"port" env:"HTTP_PORT"`
	Debug bool   `yaml:"debug" env:"DEBUG"`
}

func LoadConfig() (*Config, error) {
	cfg := &Config{}
	err := cleanenv.ReadEnv(cfg)
	if err != nil {
		return nil, errors.Wrap(err, "could not read config")
	}
	return cfg, nil
}
