package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/minguu42/mtasks/pkg/route"
)

func main() {
	r := chi.NewRouter()

	route.Route(r)

	s := &http.Server{
		Addr:           ":8080",
		Handler:        r,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

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
