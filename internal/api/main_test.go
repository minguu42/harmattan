package api_test

import (
	"context"
	"log"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/minguu42/harmattan/internal/api"
	"github.com/minguu42/harmattan/internal/atel"
	"github.com/minguu42/harmattan/internal/auth"
	"github.com/minguu42/harmattan/internal/database"
	"github.com/minguu42/harmattan/internal/database/databasetest"
	"github.com/minguu42/harmattan/internal/lib/clock"
)

const (
	testUserID = "USER-000000000000000000001"
	// token はクレームが以下の値のテスト用IDトークン
	// sub = "USER-000000000000000000001", exp = "2025-01-01 01:00:00 JST", iat = "2025-01-01 00:00:00 JST"
	token = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiJVU0VSLTAwMDAwMDAwMDAwMDAwMDAwMDAwMSIsImV4cCI6MTczNTY2MDgwMCwiaWF0IjoxNzM1NjU3MjAwfQ.Y2TZhCwHr6OosG7YM3nKObz6mDD0k6EpVrxELF7eFi8"
)

var (
	jst      = time.FixedZone("JST", 9*60*60)
	fixedNow = time.Date(2025, 1, 1, 0, 10, 0, 0, jst)

	th  http.Handler
	tdb *databasetest.ClientWithContainer
)

func init() {
	time.Local = jst
	atel.SetLogger(atel.New(os.Stdout, atel.LevelError, false))
}

func TestMain(m *testing.M) {
	ctx := context.Background()

	var err error
	tdb, err = databasetest.NewClientWithContainer(ctx, "harmattan_test")
	if err != nil {
		log.Fatalf("%+v", err)
	}
	defer atel.Capture(ctx, "Failed to close test database client")(tdb.Close)

	err = tdb.TruncateAndInsert(ctx, []any{database.Users{
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

	authn, err := auth.NewAuthenticator("cIZ15duBB4CjZNxD6CH8jBgc5sP5Ch7G", 1*time.Hour)
	if err != nil {
		log.Fatalf("%+v", err)
	}

	db, err := database.NewClient(ctx, &database.Config{
		DSN: database.DSN{
			Host:     tdb.Host,
			Port:     tdb.Port,
			Database: tdb.Database,
			User:     tdb.User,
			Password: tdb.Password,
		},
		MaxOpenConns:    25,
		MaxIdleConns:    25,
		ConnMaxLifetime: 5 * time.Minute,
	})
	if err != nil {
		log.Fatalf("%+v", err)
	}
	defer atel.Capture(ctx, "Failed to close database client")(db.Close)

	h, err := api.NewHandler(&api.Factory{
		Auth: authn,
		DB:   db,
	}, "xxxxxxx", []string{"*"})
	if err != nil {
		log.Fatalf("%+v", err)
	}
	th = fixNow(h, fixedNow)

	m.Run()
}

func fixNow(next http.Handler, tm time.Time) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, r.WithContext(clock.WithFixedNow(r.Context(), tm)))
	})
}
