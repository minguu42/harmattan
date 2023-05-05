package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

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
			shutdownErr <- fmt.Errorf("s.Shutdown failed: %w", err)
			return
		}
		shutdownErr <- nil
	}()

	if err := s.ListenAndServe(); err != http.ErrServerClosed {
		log.Fatalf("s.ListenAndServe failed: %v", err)
	}
	if err := <-shutdownErr; err != nil {
		log.Fatal(err)
	}
}
