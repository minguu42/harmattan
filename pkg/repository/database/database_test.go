package database

import (
	"context"
	"os"
	"testing"

	"github.com/minguu42/opepe/pkg/logging"
)

var testDB *DB

func TestMain(m *testing.M) {
	var err error
	testDB, err = Open(DSN("root", "", "localhost", 3306, "opepe_test"))
	if err != nil {
		logging.Errorf(context.Background(), "Open failed: %s", err)
		os.Exit(1)
	}
	defer testDB.Close()

	m.Run()
}
