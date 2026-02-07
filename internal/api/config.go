package api

import (
	"time"

	"github.com/minguu42/harmattan/internal/atel"
	"github.com/minguu42/harmattan/internal/auth"
	"github.com/minguu42/harmattan/internal/database"
)

type Config struct {
	Host           string        `env:"API_HOST" default:"0.0.0.0"`
	Port           int           `env:"API_PORT" default:"8080"`
	ReadTimeout    time.Duration `env:"API_READ_TIMEOUT" default:"2s"`
	WriteTimeout   time.Duration `env:"API_WRITE_TIMEOUT" default:"2s"`
	StopTimeout    time.Duration `env:"API_STOP_TIMEOUT" default:"25s"`
	AllowedOrigins []string      `env:"API_ALLOWED_ORIGINS" default:"http://localhost:5173,http://127.0.0.1:5173"`

	Auth auth.Config
	DB   database.Config

	LogLevel           atel.Level `env:"LOG_LEVEL" default:"info"` // "debug" || "info" || "warn" || "error"
	LogPrettyPrint     bool       `env:"LOG_PRETTY_PRINT" default:"false"`
	TraceExporter      string     `env:"TRACE_EXPORTER"` // "otlp" || "stdout" || ""
	TraceCollectorHost string     `env:"TRACE_COLLECTOR_HOST"`
	TraceCollectorPort int        `env:"TRACE_COLLECTOR_PORT"`
}
