package database

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net"
	"strconv"
	"time"

	"github.com/avast/retry-go/v4"
	"github.com/minguu42/harmattan/internal/atel"
	"github.com/minguu42/harmattan/internal/lib/errtrace"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/plugin/opentelemetry/tracing"
)

var ErrNotFound = errors.New("model not found")

type Client struct {
	gormDB *gorm.DB
}

type DSN struct {
	Host     string
	Port     int
	Database string
	User     string
	Password string
}

func (d DSN) String() string {
	return fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&loc=Local&parseTime=True",
		d.User, d.Password, net.JoinHostPort(d.Host, strconv.Itoa(d.Port)), d.Database)
}

type Config struct {
	DSN             DSN
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime time.Duration
}

func NewClient(ctx context.Context, conf *Config) (*Client, error) {
	db, err := sql.Open("mysql", conf.DSN.String())
	if err != nil {
		return nil, errtrace.Wrap(err)
	}
	db.SetMaxOpenConns(conf.MaxOpenConns)
	db.SetMaxIdleConns(conf.MaxIdleConns)
	db.SetConnMaxLifetime(conf.ConnMaxLifetime)

	ping := func() error { return db.PingContext(ctx) }
	if err := retry.Do(ping, retry.Attempts(10), retry.Context(ctx)); err != nil {
		return nil, errtrace.Wrap(err)
	}

	gormDB, err := gorm.Open(mysql.New(mysql.Config{Conn: db}), &gorm.Config{
		Logger:         customLogger{},
		TranslateError: true,
	})
	if err != nil {
		return nil, errtrace.Wrap(err)
	}
	if err := gormDB.Use(tracing.NewPlugin(tracing.WithoutMetrics(), tracing.WithoutQueryVariables())); err != nil {
		return nil, errtrace.Wrap(err)
	}
	return &Client{gormDB: gormDB}, nil
}

type customLogger struct{}

func (l customLogger) LogMode(_ logger.LogLevel) logger.Interface  { return l }
func (l customLogger) Info(_ context.Context, _ string, _ ...any)  {}
func (l customLogger) Warn(_ context.Context, _ string, _ ...any)  {}
func (l customLogger) Error(_ context.Context, _ string, _ ...any) {}
func (l customLogger) Trace(ctx context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), _ error) {
	atel.SQLLog(ctx, begin, fc)
}

func (c *Client) Close() error {
	db, err := c.gormDB.DB()
	if err != nil {
		return errtrace.Wrap(err)
	}
	return errtrace.Wrap(db.Close())
}

func (c *Client) Ping(ctx context.Context) error {
	db, err := c.gormDB.DB()
	if err != nil {
		return errtrace.Wrap(err)
	}
	return errtrace.Wrap(db.PingContext(ctx))
}

func (c *Client) db(ctx context.Context) *gorm.DB {
	return c.gormDB.WithContext(ctx)
}
