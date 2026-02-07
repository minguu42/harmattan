package main

import (
	"context"
	"net"
	"net/http"
	"os"
	"os/signal"
	"runtime/debug"
	"slices"
	"strconv"
	"syscall"
	"time"

	"github.com/minguu42/harmattan/internal/alog"
	"github.com/minguu42/harmattan/internal/api"
	"github.com/minguu42/harmattan/internal/lib/env"
	"github.com/minguu42/harmattan/internal/lib/errtrace"
)

var revision = "unknown"

func init() {
	time.Local = time.FixedZone("JST", 9*60*60)

	if info, ok := debug.ReadBuildInfo(); ok {
		if i := slices.IndexFunc(info.Settings, func(s debug.BuildSetting) bool {
			return s.Key == "vcs.revision"
		}); i != -1 {
			revision = info.Settings[i].Value[:len(revision)]
		}
	}
}

func main() {
	ctx := context.Background()
	if err := mainRun(ctx); err != nil {
		alog.Fatal(ctx, "Failed to execute mainRun", err)
	}
}

func mainRun(ctx context.Context) error {
	conf, err := env.Load[api.Config]()
	if err != nil {
		return errtrace.Wrap(err)
	}

	factory, err := api.NewFactory(ctx, conf)
	if err != nil {
		return errtrace.Wrap(err)
	}
	defer alog.Capture(ctx, "Failed to close factory")(factory.Close)

	handler, err := api.NewHandler(factory, revision, conf.AllowedOrigins)
	if err != nil {
		return errtrace.Wrap(err)
	}
	server := &http.Server{
		Addr:         net.JoinHostPort(conf.Host, strconv.Itoa(conf.Port)),
		Handler:      handler,
		ReadTimeout:  conf.ReadTimeout,
		WriteTimeout: conf.WriteTimeout,
	}

	serveErr := make(chan error)
	go func() {
		alog.Event(ctx, "Start accepting requests")
		serveErr <- server.ListenAndServe()
	}()

	sigterm := make(chan os.Signal, 1)
	signal.Notify(sigterm, syscall.SIGTERM)
	select {
	case err := <-serveErr:
		return errtrace.Wrap(err)
	case <-sigterm:
	}

	ctx, cancel := context.WithTimeout(ctx, conf.StopTimeout)
	defer cancel()

	alog.Event(ctx, "Stop accepting requests")
	if err := server.Shutdown(ctx); err != nil {
		return errtrace.Wrap(err)
	}
	alog.Event(ctx, "Server shutdown completed")
	return nil
}
