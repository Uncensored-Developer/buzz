package server

import (
	"context"
	"errors"
	"fmt"
	"github.com/Uncensored-Developer/buzz/internal/config"
	"go.uber.org/zap"
	"net"
	"net/http"
	"os"
	"sync"
	"time"
)

type Server struct {
	logger *zap.Logger
	config *config.Config
}

func NewServer(cfg *config.Config, logger *zap.Logger) *Server {
	return &Server{config: cfg, logger: logger}
}

func (s *Server) setupHandler() http.Handler {
	mux := http.NewServeMux()
	var handler http.Handler = mux

	// Middleware
	return handler
}

func (s *Server) Run(ctx context.Context) error {
	srv := s.setupHandler()

	httpServer := &http.Server{
		Addr:    net.JoinHostPort(s.config.Host, s.config.Port),
		Handler: srv,
	}

	go func() {
		logMsg := fmt.Sprintf("Listening on %s", httpServer.Addr)
		s.logger.Info(logMsg)
		//log.Printf("Listening on %s\n", httpServer.Addr)
		if err := httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			fmt.Fprintf(os.Stderr, "error listening and serving: %s\n", err)
		}
	}()

	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()
		<-ctx.Done()

		shutdownCtx := context.Background()
		shutdownCtx, cancel := context.WithTimeout(shutdownCtx, 10*time.Second)
		defer cancel()

		if err := httpServer.Shutdown(shutdownCtx); err != nil {
			fmt.Fprintf(os.Stderr, "error shutting down HTTP server: %s\n", err)
		}
	}()
	wg.Wait()
	return nil
}
