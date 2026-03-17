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
	if conf.MaxOpenConns != 0 {
		db.SetMaxOpenConns(conf.MaxOpenConns)
	}
	if conf.MaxIdleConns != 0 {
		db.SetMaxIdleConns(conf.MaxIdleConns)
	}
	if conf.ConnMaxLifetime != 0 {
		db.SetConnMaxLifetime(conf.ConnMaxLifetime)
	}

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

type txKey struct{}

func (c *Client) db(ctx context.Context) *gorm.DB {
	if tx, ok := ctx.Value(txKey{}).(*gorm.DB); ok {
		return tx.WithContext(ctx)
	}
	return c.gormDB.WithContext(ctx)
}

// Begin はトランザクションを開始する
// 戻り値の関数は *error を受け取り *error の値が nil の場合はコミット、そうでない場合はロールバックを行う
// すでにトランザクションが開始されている場合は外側のトランザクションを再利用し、部分ロールバックは行わない
func (c *Client) Begin(ctx context.Context) (context.Context, func(*error), error) {
	if _, ok := ctx.Value(txKey{}).(*gorm.DB); ok {
		return ctx, func(_ *error) {}, nil
	}

	tx := c.db(ctx).Begin()
	if err := tx.Error; err != nil {
		return ctx, nil, errtrace.Wrap(err)
	}

	return context.WithValue(ctx, txKey{}, tx), func(errp *error) {
		if errp == nil {
			panic("database: commitOrRollback called with nil error pointer")
		}

		if *errp != nil {
			// ロールバックが失敗するのは接続が切断された場合であり、その場合はDB側でロールバックされるためエラーは無視する
			tx.Rollback()
			return
		}
		if err := tx.Commit().Error; err != nil {
			*errp = errtrace.Wrap(err)
		}
	}, nil
}
