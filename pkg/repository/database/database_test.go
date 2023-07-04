package database

import (
	"os"
	"testing"

	"github.com/minguu42/mtasks/pkg/logging"
)

var testDB *DB

func TestMain(m *testing.M) {
	var err error
	testDB, err = Open(DSN("root", "", "localhost", 3306, "mtasks_test"))
	if err != nil {
		logging.Errorf("Open failed: %s", err)
		os.Exit(1)
	}
	defer testDB.Close()

	m.Run()
}
