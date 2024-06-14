package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/pkg/errors"
	"log"
	"os"
)

type Config struct {
	Host string `yaml:"host" env:"HTTP_HOST"`
	Port string `yaml:"port" env:"HTTP_PORT"`
}

func LoadConfig() (*Config, error) {
	cfg := &Config{}

	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	err = cleanenv.ReadConfig(dir+"/config.yaml", cfg)
	if err != nil {
		return nil, errors.Wrap(err, "could not read config")
	}
	return cfg, nil
}
