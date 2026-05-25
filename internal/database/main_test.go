package database_test

import (
	"context"
	"log"
	"log/slog"
	"os"
	"testing"
	"time"

	"github.com/minguu42/harmattan/internal/atel"
	"github.com/minguu42/harmattan/internal/database"
	"github.com/minguu42/harmattan/internal/database/databasetest"
)

var (
	jst *time.Location

	c   *database.Client
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
	tdb, err = databasetest.NewClient(ctx, "database_test")
	if err != nil {
		log.Fatalf("%+v", err)
	}
	defer atel.Capture(ctx, "Failed to close test database client")(tdb.Close)

	c, err = database.NewClient(ctx, &database.Config{DSN: tdb.DSN})
	if err != nil {
		log.Fatalf("%+v", err)
	}
	defer atel.Capture(ctx, "Failed to close database client")(c.Close)

	m.Run()
}
