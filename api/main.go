package main

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/minguu42/harmattan/api/factory"
	"github.com/minguu42/harmattan/api/handler"
	"github.com/minguu42/harmattan/lib/applog"
	"github.com/minguu42/harmattan/lib/env"
)

//go:generate go tool ogen -clean -config ../.ogen.yaml -package openapi -target ../internal/openapi openapi.yaml

func init() {
	time.Local = time.FixedZone("JST", 9*60*60)
}

func main() {
	ctx := context.Background()

	l := applog.New()
	if err := mainRun(ctx, l); err != nil {
		l.Error(ctx, err.Error())
		os.Exit(1)
	}
}

func mainRun(ctx context.Context, logger *applog.Logger) error {
	var conf factory.Config
	if err := env.Load(&conf); err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	f, err := factory.New(conf)
	if err != nil {
		return fmt.Errorf("failed to create factory: %w", err)
	}
	defer f.Close()

	h, err := handler.New(f, logger)
	if err != nil {
		return fmt.Errorf("failed to create handler: %w", err)
	}
	s := &http.Server{
		Addr:         net.JoinHostPort(conf.API.Host, strconv.Itoa(conf.API.Port)),
		Handler:      h,
		ReadTimeout:  conf.API.ReadTimeout,
		WriteTimeout: conf.API.WriteTimeout,
	}

	serveErr := make(chan error)
	go func() {
		logger.Event(ctx, "Start accepting requests")
		serveErr <- s.ListenAndServe()
	}()

	sigterm := make(chan os.Signal, 1)
	signal.Notify(sigterm, syscall.SIGTERM)
	select {
	case err := <-serveErr:
		return fmt.Errorf("failed to listen and serve: %w", err)
	case <-sigterm:
	}

	ctxWithTimeout, cancel := context.WithTimeout(ctx, conf.API.StopTimeout)
	defer cancel()

	logger.Event(ctx, "Begin graceful shutdown")
	if err := s.Shutdown(ctxWithTimeout); err != nil {
		return fmt.Errorf("failed to shutdown server: %w", err)
	}
	logger.Event(ctx, "Stop accepting requests")
	return nil
}
