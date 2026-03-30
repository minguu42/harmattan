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
	"github.com/stretchr/testify/require"
)

var (
	jst = time.FixedZone("JST", 9*60*60)
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

	m.Run()
}

func setupTest(t *testing.T, tableRows []any) (*database.Client, *databasetest.TestDB) {
	t.Helper()
	testDB := tdb.NewTestDB(t, tableRows)
	c, err := database.NewClient(t.Context(), database.Config{
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
	require.NoError(t, err)
	t.Cleanup(func() { c.Close() })
	return c, testDB
}
