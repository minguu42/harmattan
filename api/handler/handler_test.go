package handler_test

import (
	"context"
	"log"
	"net/http"
	"path/filepath"
	"runtime"
	"testing"
	"time"

	"github.com/ikawaha/httpcheck"
	"github.com/minguu42/harmattan/api/factory"
	"github.com/minguu42/harmattan/api/handler"
	"github.com/minguu42/harmattan/internal/auth"
	"github.com/minguu42/harmattan/internal/database"
	"github.com/minguu42/harmattan/internal/oapi"
	"github.com/minguu42/harmattan/lib/applog"
	"github.com/minguu42/harmattan/lib/clock/clocktest"
	"github.com/minguu42/harmattan/lib/databasetest"
	"github.com/minguu42/harmattan/lib/idgen/idgentest"
)

var (
	h   http.Handler
	now time.Time
	tdb *databasetest.ClientWithContainer
)

func init() {
	time.Local = time.FixedZone("JST", 9*60*60)
	now = time.Date(2025, 1, 1, 0, 10, 0, 0, time.Local)
}

func TestMain(m *testing.M) {
	ctx := context.Background()

	var err error
	tdb, err = databasetest.NewClientWithContainer(ctx, "maindb_test")
	if err != nil {
		log.Fatalf("failed to create mysql client: %s", err)
	}
	defer func() {
		if err := tdb.Close(ctx); err != nil {
			log.Fatalf("failed to close database test client: %s", err)
		}
	}()

	_, f, _, _ := runtime.Caller(0)
	if err := tdb.Migrate(ctx, filepath.Join(filepath.Dir(f), "..", "..", "infra", "mysql", "schema.sql")); err != nil {
		log.Fatalf("failed to migrate test db: %s", err)
	}

	authn, err := auth.NewAuthenticator(auth.Config{
		IDTokenExpiration: 1 * time.Hour,
		IDTokenSecret:     "cIZ15duBB4CjZNxD6CH8jBgc5sP5Ch7G",
	})
	if err != nil {
		log.Fatalf("failed to create authenticator: %s", err)
	}

	db, err := database.NewClient(database.Config{
		Host:            tdb.Host,
		Port:            tdb.Port,
		Database:        tdb.Database,
		User:            tdb.User,
		Password:        tdb.Password,
		MaxOpenConns:    25,
		MaxIdleConns:    25,
		ConnMaxLifetime: 5 * time.Minute,
	})
	if err != nil {
		log.Fatalf("failed to create database client: %s", err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			log.Fatalf("failed to close database client: %s", err)
		}
	}()

	h, err = handler.New(&factory.Factory{
		Auth: authn,
		DB:   db,
	}, applog.New())
	if err != nil {
		log.Fatalf("failed to create handler: %s", err)
	}

	m.Run()
}

func fixTimeMiddleware(next http.Handler, tm time.Time) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		ctx = clocktest.WithFixedNow(ctx, tm)
		ctx = idgentest.WithFixedULID(ctx, "01JGFJJZ000000000000000000")
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func TestHandler_NotFound(t *testing.T) {
	httpcheck.New(h).Test(t, "GET", "/non-existent-path").
		Check().
		HasStatus(404).
		HasJSON(oapi.Error{Code: 404, Message: "指定したパスは見つかりません"})
}

func TestHandler_MethodNotFound(t *testing.T) {
	t.Run("method not allowed", func(t *testing.T) {
		httpcheck.New(h).Test(t, "POST", "/health").
			Check().
			HasStatus(405).
			HasJSON(oapi.Error{Code: 405, Message: "指定したメソッドは許可されていません"})
	})
	t.Run("options", func(t *testing.T) {
		httpcheck.New(h).Test(t, "OPTIONS", "/health").
			Check().
			HasStatus(204).
			HasHeader("Access-Control-Allow-Methods", "GET").
			HasHeader("Access-Control-Allow-Headers", "Content-Type")
	})
}
