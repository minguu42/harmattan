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
	"github.com/sebdah/goldie/v2"
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

type test struct {
	id            string             // テストID（goldenファイル名）
	method        string             // リクエストメソッド
	path          string             // リクエストのURLパス
	body          io.Reader          // リクエストボディ
	statusCode    int                // ステータスコード
	prepareMockFn func(t *testing.T) // モック関数を準備する
	needsRollback bool               // この値がtrueの場合はテスト時にトランザクションをかけて、テスト後にロールバックを行う
}

func run(t *testing.T, tests []test) {
	for _, tt := range tests {
		t.Run(tt.id, func(t *testing.T) {
			if tt.needsRollback {
				if err := tdb.Begin(); err != nil {
					t.Fatalf("tdb.Begin failed: %s", err)
				}
				defer tdb.Rollback()
			}

			if tt.prepareMockFn != nil {
				tt.prepareMockFn(t)
			}

			var respBody any
			statusCode, err := doRequest(tt.method, tt.path, tt.body, &respBody)
			if err != nil {
				t.Fatalf("doRequest failed: %s", err)
			}

			if tt.statusCode != statusCode {
				t.Fatalf("status code want %d, but %d", tt.statusCode, statusCode)
			}

			if statusCode == http.StatusNoContent ||
				statusCode < 200 || 300 <= statusCode {
				return
			}

			g := goldie.New(t,
				goldie.WithFixtureDir("../../testdata"),
				goldie.WithNameSuffix(".golden.json"))
			g.AssertJson(t, tt.id, respBody)
		})
	}
}

func doRequest(method, path string, body io.Reader, respBody any) (statusCode int, err error) {
	r, err := http.NewRequest(method, ts.URL+path, body)
	if err != nil {
		return 0, fmt.Errorf("http.NewRequest failed: %w", err)
	}
	r.Header.Set("Content-Type", "application/json")
	r.Header.Set("X-Api-Key", "rAM9Fm9huuWEKLdCwHBcju9Ty_-TL2tDsAicmMrXmUnaCGp3RtywzYpMDPdEtYtR")

	c := http.Client{}
	resp, err := c.Do(r)
	if err != nil {
		return 0, fmt.Errorf("c.Do failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNoContent ||
		resp.StatusCode < 200 || 300 <= resp.StatusCode {
		return resp.StatusCode, nil
	}

	raw, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, fmt.Errorf("io.ReadAll failed: %w", err)
	}
	if err := json.Unmarshal(raw, respBody); err != nil {
		return 0, fmt.Errorf("json.Unmarshal failed: %w", err)
	}
	return resp.StatusCode, nil
}
