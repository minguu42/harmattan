package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/minguu42/mtasks/api"
)

func main() {
	db, err := api.OpenDB(api.DSN("root", "", "mtasks-db-local", 3306, "db_local"))
	if err != nil {
		api.Fatalf("rdb.New failed: %v", err)
	}

	s := api.New(db)
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

	api.Infof("Start accepting requests")
	if err := s.ListenAndServe(); err != http.ErrServerClosed {
		api.Fatalf("s.ListenAndServe failed: %v", err)
	}

	if err := <-shutdownErr; err != nil {
		api.Fatalf("s.Shutdown failed: %v", err)
	}
	api.Infof("Stop accepting requests")
}
