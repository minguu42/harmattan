package main

import (
	"cmp"
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/minguu42/harmattan/internal/alog"
	"github.com/minguu42/harmattan/internal/api/handler"
	"github.com/minguu42/harmattan/internal/factory"
	"github.com/minguu42/harmattan/internal/lib/env"
	"github.com/minguu42/harmattan/internal/lib/errcapture"
)

func init() {
	time.Local = time.FixedZone("JST", 9*60*60)
}

func main() {
	ctx := context.Background()

	level := alog.Level(cmp.Or(os.Getenv("LOG_LEVEL"), "info"))
	indent := os.Getenv("LOG_INDENT") == "true"
	l := alog.New(level, indent)
	if err := mainRun(ctx, l); err != nil {
		l.Error(ctx, err.Error())
		os.Exit(1)
	}
}

func mainRun(ctx context.Context, logger *alog.Logger) (err error) {
	var conf factory.Config
	if err := env.Load(&conf); err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	f, err := factory.New(ctx, conf, logger)
	if err != nil {
		return fmt.Errorf("failed to create factory: %w", err)
	}
	defer errcapture.Do(&err, f.Close)

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

	ctx, cancel := context.WithTimeout(ctx, conf.API.StopTimeout)
	defer cancel()

	logger.Event(ctx, "Stop accepting requests")
	if err := s.Shutdown(ctx); err != nil {
		return fmt.Errorf("failed to shutdown server: %w", err)
	}
	logger.Event(ctx, "Server shutdown completed")
	return nil
}
