package databasetest

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net"
	"os"
	"reflect"
	"strconv"
	"strings"
	"testing"

	"github.com/avast/retry-go/v4"
	"github.com/stretchr/testify/assert"
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
			Image: "mysql:8.0.30",
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
		return nil, fmt.Errorf("failed to create mysql container: %w", err)
	}

	portNet, err := mysqlC.MappedPort(ctx, "3306/tcp")
	if err != nil {
		return nil, fmt.Errorf("failed to get mapped port: %w", err)
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
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	ping := func() error { return db.PingContext(ctx) }
	if err := retry.Do(ping, retry.Attempts(10), retry.Context(ctx)); err != nil {
		return nil, fmt.Errorf("failed to connect database: %w", err)
	}

	gormDB, err := gorm.Open(mysql.New(mysql.Config{Conn: db}), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to create gorm client: %w", err)
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

func (c *ClientWithContainer) Close(ctx context.Context) error {
	if err := c.db.Close(); err != nil {
		return fmt.Errorf("failed to close db connection: %w", err)
	}
	if err := c.container.Terminate(ctx); err != nil {
		return fmt.Errorf("failed to terminate database: %w", err)
	}
	return nil
}

func (c *ClientWithContainer) Migrate(ctx context.Context, name string) error {
	data, err := os.ReadFile(name)
	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}

	for query := range strings.SplitSeq(string(data), ";") {
		if query = strings.TrimSpace(query); query == "" {
			continue
		}
		if !strings.HasPrefix(strings.ToUpper(query), "CREATE TABLE") {
			return errors.New("only CREATE TABLE statements are supported")
		}

		if _, err := c.db.ExecContext(ctx, query); err != nil {
			return fmt.Errorf("failed to execute query: %w", err)
		}
	}
	return nil
}

func (c *ClientWithContainer) Insert(ctx context.Context, data []any) error {
	for _, rows := range data {
		if err := c.gormDB.WithContext(ctx).Create(rows).Error; err != nil {
			return fmt.Errorf("failed to insert rows: %w", err)
		}
	}
	return nil
}

func (c *ClientWithContainer) Reset(ctx context.Context, data []any) error {
	for _, table := range data {
		// 何の条件もなしに一括削除を行えないため、"WHERE 1 = 1"で回避している
		if err := c.gormDB.WithContext(ctx).Where("1 = 1").Delete(table).Error; err != nil {
			return fmt.Errorf("failed to delete rows: %w", err)
		}
	}
	return nil
}

func (c *ClientWithContainer) Assert(t *testing.T, ctx context.Context, data []any) {
	for _, want := range data {
		rv := reflect.ValueOf(want)
		if rv.Kind() != reflect.Slice {
			log.Fatalf("expecting slice type, got %s", rv.Kind())
		}

		elemType := rv.Type().Elem()
		slicePtr := reflect.New(reflect.SliceOf(elemType)).Interface()
		if err := c.gormDB.WithContext(ctx).Find(slicePtr).Error; err != nil {
			log.Fatal(err)
		}

		actualSlice := reflect.ValueOf(slicePtr).Elem()
		if actualSlice.Len() != rv.Len() {
			log.Fatalf("expecting %d elements, got %d", rv.Len(), actualSlice.Len())
		}

		got := actualSlice.Interface()
		assert.Equal(t, want, got)
	}
}
