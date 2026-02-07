package database_test

import (
	"context"
	"log"
	"os"
	"testing"
	"time"

	"github.com/minguu42/harmattan/internal/alog"
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
	alog.SetLogger(alog.New(os.Stdout, alog.LevelError, false))
}

func TestMain(m *testing.M) {
	ctx := context.Background()

	var err error
	tdb, err = databasetest.NewClientWithContainer(ctx, "database_test")
	if err != nil {
		log.Fatalf("%+v", err)
	}
	defer alog.Capture(ctx, "Failed to close test database client")(tdb.Close)

	c, err = database.NewClient(ctx, database.Config{
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
	defer alog.Capture(ctx, "Failed to close database client")(c.Close)

	m.Run()
}
