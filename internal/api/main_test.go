package api_test

import (
	"context"
	"log"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/minguu42/harmattan/internal/api"
	"github.com/minguu42/harmattan/internal/atel"
	"github.com/minguu42/harmattan/internal/database/databasetest"
	"github.com/minguu42/harmattan/internal/lib/clock"
	"github.com/minguu42/harmattan/internal/lib/idgen"
)

const (
	fixedID = "GENERATED-ID-0000000000001"
	// token はクレームが以下の値のテスト用IDトークン
	// sub = "USER-000000000000000000001", exp = "2025-01-01 01:00:00 JST", iat = "2025-01-01 00:00:00 JST"
	token = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiJVU0VSLTAwMDAwMDAwMDAwMDAwMDAwMDAwMSIsImV4cCI6MTczNTY2MDgwMCwiaWF0IjoxNzM1NjU3MjAwfQ.Y2TZhCwHr6OosG7YM3nKObz6mDD0k6EpVrxELF7eFi8"
)

var (
	jst *time.Location
	ts  *httptest.Server
	tdb *databasetest.Client
)

func init() {
	var err error
	jst, err = time.LoadLocation("Asia/Tokyo")
	if err != nil {
		log.Fatalf("failed to load location: %v", err)
	}
	time.Local = jst
	atel.SetLogger(atel.New(os.Stdout, slog.LevelError, false))
}

func TestMain(m *testing.M) {
	ctx := context.Background()

	var err error
	tdb, err = databasetest.NewClient(ctx, "harmattan_test")
	if err != nil {
		log.Fatalf("%+v", err)
	}
	defer atel.Capture(ctx, "Failed to close test database client")(tdb.Close)

	f, err := api.NewFactory(ctx, &api.Config{
		IDTokenSecret:     "cIZ15duBB4CjZNxD6CH8jBgc5sP5Ch7G",
		IDTokenExpiration: 1 * time.Hour,
		DBHost:            tdb.DSN.Host,
		DBPort:            tdb.DSN.Port,
		DBDatabase:        tdb.DSN.Database,
		DBUser:            tdb.DSN.User,
		DBPassword:        tdb.DSN.Password,
	})
	if err != nil {
		log.Fatalf("%+v", err)
	}
	defer atel.Capture(ctx, "Failed to close factory")(f.Close)

	h, err := api.NewHandler(f, "xxxxxxx", []string{"*"})
	if err != nil {
		log.Fatalf("%+v", err)
	}

	ts = httptest.NewServer(fixNow(fixID(h, fixedID), time.Date(2025, 1, 1, 0, 10, 0, 0, jst)))
	defer ts.Close()

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
