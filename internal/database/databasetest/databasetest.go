package databasetest

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"strconv"
	"strings"
	"testing"

	"github.com/avast/retry-go/v4"
	"github.com/minguu42/harmattan/internal/lib/errtrace"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type ClientWithContainer struct {
	container testcontainers.Container
	db        *sql.DB
	gormDB    *gorm.DB

	Host     string
	Port     int
	Database string
	User     string
	Password string
}

func NewClientWithContainer(ctx context.Context, database string) (*ClientWithContainer, error) {
	mysqlC, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			Image: "mysql:8.0.42",
			Env: map[string]string{
				"MYSQL_DATABASE":             database,
				"MYSQL_ALLOW_EMPTY_PASSWORD": "yes",
			},
			ExposedPorts: []string{"3306/tcp"},
			WaitingFor:   wait.ForListeningPort("3306/tcp"),
		},
		Started: true,
	})
	if err != nil {
		return nil, errtrace.Wrap(err)
	}

	portNet, err := mysqlC.MappedPort(ctx, "3306/tcp")
	if err != nil {
		return nil, errtrace.Wrap(err)
	}

	host := "127.0.0.1"
	port := portNet.Int()
	user := "root"
	password := ""
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?parseTime=true&charset=utf8mb4&loc=Local",
		user,
		password,
		net.JoinHostPort(host, strconv.Itoa(port)),
		database,
	)
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, errtrace.Wrap(err)
	}

	ping := func() error { return db.PingContext(ctx) }
	if err := retry.Do(ping, retry.Attempts(10), retry.Context(ctx)); err != nil {
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
	return &ClientWithContainer{
		container: mysqlC,
		db:        db,
		gormDB:    gormDB,
		Host:      host,
		Port:      port,
		Database:  database,
		User:      user,
		Password:  password,
	}, nil
}

func (c *ClientWithContainer) Close() error {
	if err := c.db.Close(); err != nil {
		return errtrace.Wrap(err)
	}
	if err := c.container.Terminate(context.Background()); err != nil {
		return errtrace.Wrap(err)
	}
	return nil
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

func (c *ClientWithContainer) TruncateAndInsert(ctx context.Context, tableRows []any) error {
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

func (c *ClientWithContainer) Assert(t *testing.T, data []any) {
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
