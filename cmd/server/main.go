package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/minguu42/mtasks/pkg/app"
	"github.com/minguu42/mtasks/pkg/database"
	"github.com/minguu42/mtasks/pkg/env"
	"github.com/minguu42/mtasks/pkg/logging"
)

func main() {
	appEnv, err := env.Load()
	if err != nil {
		logging.Fatalf("env.Load failed: %v", err)
	}

	db, err := database.Open(context.Background(), appEnv.MySQL.DSN())
	if err != nil {
		logging.Fatalf("database.Open failed: %v", err)
	}
	defer db.Close()

	s, err := app.NewServer(appEnv.API, db)
	if err != nil {
		logging.Fatalf("server.NewServer failed: %v", err)
	}

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

	logging.Infof("Start accepting requests")
	if err := s.ListenAndServe(); err != http.ErrServerClosed {
		logging.Fatalf("s.ListenAndServe failed: %v", err)
	}

	if err := <-shutdownErr; err != nil {
		logging.Fatalf("s.Shutdown failed: %v", err)
	}
	logging.Infof("Stop accepting requests")
}
