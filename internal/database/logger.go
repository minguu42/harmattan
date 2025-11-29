package database

import (
	"context"
	"time"

	"github.com/minguu42/harmattan/internal/alog"
	"gorm.io/gorm/logger"
)

type customLogger struct{}

func (l customLogger) LogMode(_ logger.LogLevel) logger.Interface {
	return l
}

func (l customLogger) Info(_ context.Context, _ string, _ ...any) {}

func (l customLogger) Warn(_ context.Context, _ string, _ ...any) {}

func (l customLogger) Error(_ context.Context, _ string, _ ...any) {}

func (l customLogger) Trace(ctx context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), _ error) {
	alog.GORMTrace(ctx, begin, fc)
}
