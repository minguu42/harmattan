package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/minguu42/mtasks/app"
	"github.com/minguu42/mtasks/app/env"
	"github.com/minguu42/mtasks/app/logging"
	"github.com/minguu42/mtasks/app/ogen"
	"github.com/minguu42/mtasks/app/repository/database"
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

	h, err := ogen.NewServer(
		&app.Handler{Repository: db},
		&app.SecurityHandler{Repository: db},
	)
	if err != nil {
		logging.Fatalf("ogen.NewServer failed: %v", err)
	}
	s := &http.Server{
		Addr:              fmt.Sprintf("%s:%d", appEnv.API.Host, appEnv.API.Port),
		Handler:           app.LogMiddleware(h),
		ReadTimeout:       10 * time.Second,
		ReadHeaderTimeout: 10 * time.Second,
		MaxHeaderBytes:    1 << 20,
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
