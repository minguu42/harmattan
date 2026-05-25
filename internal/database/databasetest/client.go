package databasetest

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"strings"
	"testing"
	"time"

	"github.com/avast/retry-go/v5"
	"github.com/minguu42/harmattan/internal/database"
	"github.com/minguu42/harmattan/internal/lib/errtrace"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Client struct {
	container *testcontainers.DockerContainer
	db        *sql.DB
	gormDB    *gorm.DB
	DSN       database.DSN
}

func NewClient(ctx context.Context, databaseName string) (*Client, error) {
	user := "harmattan"
	password := "R2b87Yy6owa5Jxo7EkR8"
	container, err := testcontainers.Run(ctx, "mysql:8.0.42",
		testcontainers.WithEnv(map[string]string{
			"MYSQL_DATABASE":      databaseName,
			"MYSQL_USER":          user,
			"MYSQL_PASSWORD":      password,
			"MYSQL_ROOT_PASSWORD": password,
		}),
		testcontainers.WithExposedPorts("3306/tcp"),
		testcontainers.WithWaitStrategy(wait.ForListeningPort("3306/tcp")),
	)
	if err != nil {
		return nil, errtrace.Wrap(err)
	}

	portNet, err := container.MappedPort(ctx, "3306/tcp")
	if err != nil {
		return nil, errtrace.Wrap(err)
	}

	dsn := database.DSN{
		Host:     "localhost",
		Port:     int(portNet.Num()),
		Database: databaseName,
		User:     user,
		Password: password,
	}
	db, err := sql.Open("mysql", dsn.String())
	if err != nil {
		return nil, errtrace.Wrap(err)
	}

	ping := func() error { return db.PingContext(ctx) }
	if err := retry.New(retry.DelayType(retry.FixedDelay), retry.Delay(time.Second)).Do(ping); err != nil {
		return nil, errtrace.Wrap(err)
	}

	if err := applySchema(ctx, db); err != nil {
		return nil, errtrace.Wrap(err)
	}
	if _, err := db.ExecContext(ctx, "set FOREIGN_KEY_CHECKS = 0"); err != nil {
		return nil, errtrace.Wrap(err)
	}

	gormDB, err := gorm.Open(mysql.New(mysql.Config{Conn: db}), &gorm.Config{})
	if err != nil {
		return nil, errtrace.Wrap(err)
	}
	return &Client{
		container: container,
		db:        db,
		gormDB:    gormDB,
		DSN:       dsn,
	}, nil
}

func applySchema(ctx context.Context, db *sql.DB) error {
	_, f, _, _ := runtime.Caller(0)
	data, err := os.ReadFile(filepath.Join(filepath.Dir(f), "..", "..", "..", "infra", "mysql", "schema.sql"))
	if err != nil {
		return errtrace.Wrap(err)
	}

	for query := range strings.SplitSeq(string(data), ";") {
		if query = strings.TrimSpace(query); query == "" {
			continue
		}
		if !strings.HasPrefix(strings.ToUpper(query), "CREATE TABLE") {
			return errtrace.Wrap(errors.New("only CREATE TABLE statements are supported"))
		}

		if _, err := db.ExecContext(ctx, query); err != nil {
			return errtrace.Wrap(err)
		}
	}
	return nil
}

func (c *Client) Close() error {
	dbErr := c.db.Close()
	containerErr := testcontainers.TerminateContainer(c.container)
	return errtrace.Wrap(errors.Join(dbErr, containerErr))
}

func (c *Client) TruncateAndInsert(ctx context.Context, tableRows []any) error {
	for _, rows := range tableRows {
		stmt := &gorm.Statement{DB: c.gormDB}
		if err := stmt.Parse(rows); err != nil {
			return errtrace.Wrap(err)
		}
		table := stmt.Schema.Table

		if _, err := c.db.ExecContext(ctx, fmt.Sprintf("truncate table %s", table)); err != nil {
			return errtrace.Wrap(err)
		}

		if rv := reflect.ValueOf(rows); rv.Len() == 0 {
			continue
		}
		if err := c.gormDB.WithContext(ctx).Table(table).Create(rows).Error; err != nil {
			return errtrace.Wrap(err)
		}
	}
	return nil
}

func (c *Client) Assert(t *testing.T, data []any) {
	t.Helper()

	for _, want := range data {
		rv := reflect.ValueOf(want)
		if rv.Kind() != reflect.Slice {
			t.Fatalf("want slice type, got: %T", want)
		}

		gotPointer := reflect.New(reflect.SliceOf(rv.Type().Elem())).Interface()
		err := c.gormDB.Find(gotPointer).Error
		require.NoError(t, err)

		got := reflect.ValueOf(gotPointer).Elem().Interface()
		assert.ElementsMatch(t, want, got)
	}
}
