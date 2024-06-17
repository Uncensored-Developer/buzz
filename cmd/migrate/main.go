package main

import (
	"github.com/Uncensored-Developer/buzz/pkg/config"
	"github.com/Uncensored-Developer/buzz/pkg/logger"
	"github.com/Uncensored-Developer/buzz/pkg/migrate"
	"go.uber.org/zap"
	"log"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalln(err)
	}
	zapLogger := logger.NewLogger(cfg)

	err = migrate.Up(cfg.DatabaseURL, "db/migrations")
	if err != nil {
		zapLogger.Fatal("migration failed",
			zap.Error(err))
	}
	zapLogger.Info("migration successful.")
}
