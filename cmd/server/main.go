package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/minguu42/mtasks/pkg/app"
	"github.com/minguu42/mtasks/pkg/logging"
	"github.com/minguu42/mtasks/pkg/server"
)

func main() {
	if err := app.OpenDB(app.DSN("root", "", "mtasks-db-local", 3306, "db_local")); err != nil {
		logging.Fatalf("api.OpenDB failed: %v", err)
	}
	defer app.CloseDB()

	s := server.NewServer()
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
