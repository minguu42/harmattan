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
	jst = time.FixedZone("JST", 9*60*60)

	c   *database.Client
	tdb *databasetest.ClientWithContainer
)

func init() {
	time.Local = jst
	atel.SetLogger(atel.New(os.Stdout, slog.LevelError, false))
}

func TestMain(m *testing.M) {
	ctx := context.Background()

	var err error
	tdb, err = databasetest.NewClientWithContainer(ctx, "database_test")
	if err != nil {
		log.Fatalf("%+v", err)
	}
	defer atel.Capture(ctx, "Failed to close test database client")(tdb.Close)

	c, err = database.NewClient(ctx, &database.Config{
		DSN: database.DSN{
			Host:     tdb.Host,
			Port:     tdb.Port,
			Database: tdb.Database,
			User:     tdb.User,
			Password: tdb.Password,
		},
	})
	if err != nil {
		log.Fatalf("%+v", err)
	}
	defer atel.Capture(ctx, "Failed to close database client")(c.Close)

	m.Run()
}
