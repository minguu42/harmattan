package main

import (
	"context"
	"encoding/json"
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
	db.SetIDGenerator(&ulidgen.Generator{})
	if err != nil {
		logging.Fatalf(ctx, "database.Open failed: %v", err)
	}
	defer db.Close()

	h, err := ogen.NewServer(
		&handler.Handler{Repository: db},
		&handler.Security{Repository: db},
		ogen.WithNotFound(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusNotFound)
			_ = json.NewEncoder(w).Encode(ogen.Error{
				Message: "Not Found",
				Debug:   "指定したパスに対応するオペレーションは存在しない",
			})
		}),
		ogen.WithMethodNotAllowed(func(w http.ResponseWriter, r *http.Request, allowed string) {
			w.Header().Set("Content-Type", "text/plain")
			w.Header().Set("Allow", allowed)
			w.WriteHeader(http.StatusMethodNotAllowed)
			_ = json.NewEncoder(w).Encode(ogen.Error{
				Message: "Method Not Allowed",
				Debug:   fmt.Sprintf("このパスに対応しているメソッドは%sのみである", allowed),
			})
		}),
		ogen.WithErrorHandler(handler.ErrorHandler),
	)
	if err != nil {
		logging.Fatalf(ctx, "ogen.NewServer failed: %v", err)
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

	logging.Infof(ctx, "Start accepting requests")
	if err := s.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
		logging.Fatalf(ctx, "s.ListenAndServe failed: %v", err)
	}

	if err := <-shutdownErr; err != nil {
		logging.Fatalf(ctx, "s.Shutdown failed: %v", err)
	}
	logging.Infof(ctx, "Stop accepting requests")
}
