package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
)

func NewLogger() *zap.Logger {
	var zapConfig zap.Config

	if os.Getenv("BUZZ_DEBUG") == "" || os.Getenv("BUZZ_DEBUG") == "true" {
		zapConfig = zap.NewDevelopmentConfig()
	} else {
		zapConfig = zap.NewProductionConfig()
	}

	zapConfig.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	return zap.Must(zapConfig.Build())
}
