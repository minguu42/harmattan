package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/minguu42/mtasks/pkg/logger"
	"github.com/minguu42/mtasks/pkg/server"
)

func main() {
	s := server.New()

	shutdownErr := make(chan error, 1)
	go func() {
		sigterm := make(chan os.Signal, 1)
		signal.Notify(sigterm, syscall.SIGTERM)
		<-sigterm

		if err := s.Shutdown(context.Background()); err != nil {
			shutdownErr <- err
			return
		}
		shutdownErr <- nil
	}()

	if err := s.ListenAndServe(); err != http.ErrServerClosed {
		logger.Fatalf("s.ListenAndServe failed: %v\n", err)
	}
	if err := <-shutdownErr; err != nil {
		logger.Fatalf("s.Shutdown failed: %v\n", err)
	}
}
