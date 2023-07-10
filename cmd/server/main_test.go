package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/minguu42/mtasks/gen/ogen"
	"github.com/minguu42/mtasks/pkg/handler"
	"github.com/minguu42/mtasks/pkg/logging"
	"github.com/minguu42/mtasks/pkg/repository/database"
	"github.com/minguu42/mtasks/pkg/ttime"
)

var (
	ts  *httptest.Server
	tdb *database.DB
)

func TestMain(m *testing.M) {
	ctx := context.Background()
	dsn := database.DSN("root", "", "localhost", 3306, "mtasks_test")
	var err error
	tdb, err = database.Open(dsn)
	if err != nil {
		logging.Fatalf(ctx, "database.Open failed: %s", err)
	}
	defer tdb.Close()

	h, err := ogen.NewServer(
		&handler.Handler{Repository: tdb},
		&handler.Security{Repository: tdb},
	)
	if err != nil {
		logging.Fatalf(ctx, "ogen.NewServer failed: %s", err)
	}
	ts = httptest.NewServer(ttime.MiddlewareFixTime(h, time.Date(2020, 1, 2, 0, 0, 0, 0, time.UTC)))
	defer ts.Close()

	m.Run()
}

type request struct {
	method string
	path   string
	body   io.Reader
}

type response struct {
	statusCode int
	body       any
}

type test struct {
	name     string
	request  request
	response response
}

func doTestRequest(req request, got any) (*http.Response, error) {
	r, err := http.NewRequest(req.method, ts.URL+req.path, req.body)
	if err != nil {
		return nil, fmt.Errorf("http.NewRequest failed: %w", err)
	}
	r.Header.Set("Content-Type", "application/json")
	r.Header.Set("X-Api-Key", "rAM9Fm9huuWEKLdCwHBcju9Ty_-TL2tDsAicmMrXmUnaCGp3RtywzYpMDPdEtYtR")

	c := &http.Client{}
	resp, err := c.Do(r)
	if err != nil {
		return nil, fmt.Errorf("c.Do failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNoContent {
		return resp, nil
	}

	bs, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("io.ReadAll failed: %w", err)
	}
	if err := json.Unmarshal(bs, got); err != nil {
		return nil, fmt.Errorf("json.Unmarshal failed: %w", err)
	}
	return resp, nil
}
