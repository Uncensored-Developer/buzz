package logger

import (
	"github.com/Uncensored-Developer/buzz/internal/config"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func NewLogger(cfg *config.Config) *zap.Logger {
	var zapConfig zap.Config
	if cfg.Debug {
		zapConfig = zap.NewDevelopmentConfig()
	} else {
		zapConfig = zap.NewProductionConfig()
	}

	zapConfig.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	return zap.Must(zapConfig.Build())
}
