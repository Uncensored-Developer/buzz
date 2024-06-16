package server

import (
	"github.com/Uncensored-Developer/buzz/pkg/config"
	"github.com/Uncensored-Developer/buzz/pkg/logger"
	"github.com/pkg/errors"
)

func InitializeServer() (*Server, error) {
	cfg, err := config.LoadConfig()
	if err != nil {
		return nil, errors.Wrap(err, "load config failed")
	}
	zapLogger := logger.NewLogger(cfg)
	return NewServer(cfg, zapLogger), err
}
