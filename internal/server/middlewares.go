package server

import (
	"context"
	"go.uber.org/zap"
	"net/http"
	"time"
)

func RequestLoggingMiddleWare(
	ctx context.Context,
	logger *zap.Logger,
	next http.Handler,
) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			next.ServeHTTP(w, r)
			logger.Info("request completed",
				zap.String("method", r.Method),
				zap.String("url", r.URL.String()),
				zap.String("remote_addr", r.RemoteAddr),
				zap.Duration("duration", time.Since(start)),
			)
		},
	)
}
