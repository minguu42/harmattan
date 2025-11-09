package main

import (
	"cmp"
	"context"
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
	"github.com/minguu42/harmattan/internal/lib/errtrace"
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
		l.Error(ctx, "failed to execute mainRun", err)
		os.Exit(1)
	}
}

func mainRun(ctx context.Context, l *alog.Logger) error {
	var conf factory.Config
	if err := env.Load(&conf); err != nil {
		return errtrace.Wrap(err)
	}

	f, err := factory.New(ctx, conf, l)
	if err != nil {
		return errtrace.Wrap(err)
	}
	defer l.Capture(ctx, "failed to close factory")(f.Close)

	h, err := handler.New(f, l)
	if err != nil {
		return errtrace.Wrap(err)
	}
	s := &http.Server{
		Addr:         net.JoinHostPort(conf.API.Host, strconv.Itoa(conf.API.Port)),
		Handler:      h,
		ReadTimeout:  conf.API.ReadTimeout,
		WriteTimeout: conf.API.WriteTimeout,
	}

	serveErr := make(chan error)
	go func() {
		l.Event(ctx, "Start accepting requests")
		serveErr <- s.ListenAndServe()
	}()

	sigterm := make(chan os.Signal, 1)
	signal.Notify(sigterm, syscall.SIGTERM)
	select {
	case err := <-serveErr:
		return errtrace.Wrap(err)
	case <-sigterm:
	}

	ctx, cancel := context.WithTimeout(ctx, conf.API.StopTimeout)
	defer cancel()

	l.Event(ctx, "Stop accepting requests")
	if err := s.Shutdown(ctx); err != nil {
		return errtrace.Wrap(err)
	}
	l.Event(ctx, "Server shutdown completed")
	return nil
}
