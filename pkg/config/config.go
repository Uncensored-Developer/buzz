package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/pkg/errors"
	"os"
	"path/filepath"
)

type Config struct {
	Host  string `yaml:"host" env:"HTTP_HOST"`
	Port  string `yaml:"port" env:"HTTP_PORT"`
	Debug bool   `yaml:"debug" env:"DEBUG"`
}

func LoadConfig() (*Config, error) {
	cfg := &Config{}

	dir, err := os.Getwd()
	if err != nil {
		return nil, errors.Wrap(err, "could not get current dir")
	}
	parentTop := filepath.Dir(dir)

	err = cleanenv.ReadConfig(parentTop+"/config/config.yaml", cfg)
	if err != nil {
		return nil, errors.Wrap(err, "could not read config")
	}
	return cfg, nil
}
