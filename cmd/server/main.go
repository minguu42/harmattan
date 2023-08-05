package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/minguu42/opepe/gen/ogen"
	"github.com/minguu42/opepe/pkg/env"
	"github.com/minguu42/opepe/pkg/handler"
	"github.com/minguu42/opepe/pkg/handler/middleware"
	"github.com/minguu42/opepe/pkg/idgen/ulidgen"
	"github.com/minguu42/opepe/pkg/logging"
	"github.com/minguu42/opepe/pkg/repository/database"
)

func main() {
	if err := env.Load(); err != nil {
		logging.Fatalf(context.Background(), "env.Load failed: %v", err)
	}

	ctx := context.Background()
	e, err := env.Get()
	if err != nil {
		logging.Fatalf(ctx, "env.Get failed: %v", err)
	}

	dsn := database.DSN(e.MySQL.User, e.MySQL.Password, e.MySQL.Host, e.MySQL.Port, e.MySQL.Database)
	db, err := database.Open(dsn)
	if err != nil {
		logging.Fatalf(ctx, "database.Open failed: %v", err)
	}
	defer db.Close()

	h, err := ogen.NewServer(
		&handler.Handler{
			Repository:  db,
			IDGenerator: &ulidgen.Generator{},
		},
		&handler.Security{Repository: db},
		ogen.WithNotFound(handler.NotFound),
		ogen.WithMethodNotAllowed(handler.MethodNotAllowed),
		ogen.WithErrorHandler(handler.ErrorHandler),
	)
	if err != nil {
		logging.Fatalf(ctx, "ogen.NewServer failed: %v", err)
	}
	s := &http.Server{
		Addr:              fmt.Sprintf("%s:%d", e.API.Host, e.API.Port),
		Handler:           middleware.LogMiddleware(h),
		ReadTimeout:       10 * time.Second,
		ReadHeaderTimeout: 10 * time.Second,
		MaxHeaderBytes:    1 << 20,
	}

	go func() {
		logging.Infof(ctx, "Start accepting requests")
		if err := s.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			logging.Errorf(ctx, "s.ListenAndServe failed: %s", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, os.Interrupt)
	<-quit
	if err := s.Shutdown(ctx); err != nil {
		logging.Fatalf(ctx, "s.Shutdown failed: %s", err)
	}
	logging.Infof(ctx, "Stop accepting requests")
}
