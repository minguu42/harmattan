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
	"github.com/minguu42/harmattan/internal/lib/errtrace"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var ErrModelNotFound = errors.New("model not found in database")

type Client struct {
	gormDB *gorm.DB
}

type Config struct {
	Host            string        `env:"DB_HOST,required"`
	Port            int           `env:"DB_PORT,required"`
	Database        string        `env:"DB_DATABASE,required"`
	User            string        `env:"DB_USER,required"`
	Password        string        `env:"DB_PASSWORD,required"`
	MaxOpenConns    int           `env:"DB_MAX_OPEN_CONNS" default:"25"`
	MaxIdleConns    int           `env:"DB_MAX_IDLE_CONNS" default:"25"`
	ConnMaxLifetime time.Duration `env:"DB_CONN_MAX_LIFETIME" default:"5m"`
}

func NewClient(ctx context.Context, conf Config) (*Client, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&loc=Local&parseTime=True",
		conf.User,
		conf.Password,
		net.JoinHostPort(conf.Host, strconv.Itoa(conf.Port)),
		conf.Database,
	)
	db, err := sql.Open("mysql", dsn)
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
	return &Client{gormDB: gormDB}, nil
}

func (c *Client) Close() error {
	db, err := c.gormDB.DB()
	if err != nil {
		return errtrace.Wrap(err)
	}
	return errtrace.Wrap(db.Close())
}

func (c *Client) db(ctx context.Context) *gorm.DB {
	return c.gormDB.WithContext(ctx)
}
