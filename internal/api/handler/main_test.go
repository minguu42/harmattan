package handler_test

import (
	"context"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"testing"
	"time"

	"github.com/minguu42/harmattan/internal/alog"
	"github.com/minguu42/harmattan/internal/api/handler"
	"github.com/minguu42/harmattan/internal/auth"
	"github.com/minguu42/harmattan/internal/database"
	"github.com/minguu42/harmattan/internal/database/databasetest"
	"github.com/minguu42/harmattan/internal/factory"
	"github.com/minguu42/harmattan/internal/lib/clock"
	"github.com/minguu42/harmattan/internal/lib/idgen"
)

const (
	testUserID = "USER-000000000000000000001"
	// token はクレームが以下の値のテスト用IDトークン
	// sub = "USER-000000000000000000001", exp = "2025-01-01 01:00:00 JST", iat = "2025-01-01 00:00:00 JST"
	token   = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiJVU0VSLTAwMDAwMDAwMDAwMDAwMDAwMDAwMSIsImV4cCI6MTczNTY2MDgwMCwiaWF0IjoxNzM1NjU3MjAwfQ.Y2TZhCwHr6OosG7YM3nKObz6mDD0k6EpVrxELF7eFi8"
	fixedID = "GENERATED-ID-0000000000001"
)

var (
	jst      = time.FixedZone("JST", 9*60*60)
	fixedNow = time.Date(2025, 1, 1, 0, 10, 0, 0, jst)

	th  http.Handler // 現在時刻と生成されるIDは固定化されている
	tdb *databasetest.ClientWithContainer
)

func init() {
	time.Local = jst
	alog.SetDefaultLogger(alog.New(os.Stdout, alog.LevelError, false))
}

func TestMain(m *testing.M) {
	ctx := context.Background()

	var err error
	tdb, err = databasetest.NewClientWithContainer(ctx, "maindb_test")
	if err != nil {
		log.Fatalf("%+v", err)
	}
	defer alog.Capture(ctx, "Failed to close test database client")(tdb.Close)

	_, f, _, _ := runtime.Caller(0)
	if err := tdb.Migrate(ctx, filepath.Join(filepath.Dir(f), "..", "..", "..", "infra", "mysql", "schema.sql")); err != nil {
		log.Fatalf("%+v", err)
	}

	err = tdb.Insert(ctx, []any{database.Users{
		{
			ID:             testUserID,
			Email:          "user1@dummy.invalid",
			HashedPassword: "password",
			CreatedAt:      time.Date(2025, 1, 1, 0, 0, 1, 0, jst),
			UpdatedAt:      time.Date(2025, 1, 1, 0, 0, 1, 0, jst),
		},
		{
			ID:             "USER-000000000000000000002",
			Email:          "user2@dummy.invalid",
			HashedPassword: "password",
			CreatedAt:      time.Date(2025, 1, 1, 0, 0, 2, 0, jst),
			UpdatedAt:      time.Date(2025, 1, 1, 0, 0, 2, 0, jst),
		},
	}})
	if err != nil {
		log.Fatalf("%+v", err)
	}

	authn, err := auth.NewAuthenticator(auth.Config{
		IDTokenExpiration: 1 * time.Hour,
		IDTokenSecret:     "cIZ15duBB4CjZNxD6CH8jBgc5sP5Ch7G",
	})
	if err != nil {
		log.Fatalf("%+v", err)
	}

	db, err := database.NewClient(ctx, database.Config{
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
		log.Fatalf("%+v", err)
	}
	defer alog.Capture(ctx, "Failed to close database client")(db.Close)

	h, err := handler.New(&factory.Factory{
		Auth: authn,
		DB:   db,
	})
	if err != nil {
		log.Fatalf("%+v", err)
	}
	th = fixNow(fixID(h, fixedID), fixedNow)

	m.Run()
}

func fixID(next http.Handler, id string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, r.WithContext(idgen.WithFixedULID(r.Context(), id)))
	})
}

func fixNow(next http.Handler, tm time.Time) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, r.WithContext(clock.WithFixedNow(r.Context(), tm)))
	})
}
