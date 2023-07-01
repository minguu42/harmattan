package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/minguu42/mtasks/gen/ogen"
	"github.com/minguu42/mtasks/pkg/env"
	"github.com/minguu42/mtasks/pkg/handler"
	"github.com/minguu42/mtasks/pkg/logging"
	"github.com/minguu42/mtasks/pkg/repository/database"
)

func init() {
	if err := env.Load(); err != nil {
		logging.Fatalf("env.Load failed: %v", err)
	}
}

func main() {
	e, err := env.Get()
	if err != nil {
		logging.Fatalf("env.Get failed: %v", err)
	}

	dsn := database.DSN(e.MySQL.User, e.MySQL.Password, e.MySQL.Host, e.MySQL.Port, e.MySQL.Database)
	db, err := database.Open(dsn)
	if err != nil {
		logging.Fatalf("database.Open failed: %v", err)
	}
	defer db.Close()

	h, err := ogen.NewServer(
		&handler.Handler{Repository: db},
		&handler.Security{Repository: db},
	)
	if err != nil {
		logging.Fatalf("ogen.NewServer failed: %v", err)
	}
	s := &http.Server{
		Addr:              fmt.Sprintf("%s:%d", e.API.Host, e.API.Port),
		Handler:           handler.MiddlewareLog(h),
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
