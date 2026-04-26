package api

import "time"

// Config はAPIサーバの設定値を保持する構造体である
// フィールドのデフォルト値はそのまま本番環境で適用される値であり、変更は注意してください
type Config struct {
	Host           string        `env:"API_HOST" default:"0.0.0.0"`
	Port           int           `env:"API_PORT" default:"8080"`
	ReadTimeout    time.Duration `env:"API_READ_TIMEOUT" default:"2s"`
	WriteTimeout   time.Duration `env:"API_WRITE_TIMEOUT" default:"2s"`
	StopTimeout    time.Duration `env:"API_STOP_TIMEOUT" default:"25s"`
	AllowedOrigins []string      `env:"API_ALLOWED_ORIGINS,required"`

	IDTokenSecret     string        `env:"ID_TOKEN_SECRET,required"`
	IDTokenExpiration time.Duration `env:"ID_TOKEN_EXPIRATION" default:"1h"`

	DBHost            string        `env:"DB_HOST,required"`
	DBPort            int           `env:"DB_PORT,required"`
	DBDatabase        string        `env:"DB_DATABASE,required"`
	DBUser            string        `env:"DB_USER,required"`
	DBPassword        string        `env:"DB_PASSWORD,required"`
	DBMaxOpenConns    int           `env:"DB_MAX_OPEN_CONNS" default:"25"`
	DBMaxIdleConns    int           `env:"DB_MAX_IDLE_CONNS" default:"25"`
	DBConnMaxLifetime time.Duration `env:"DB_CONN_MAX_LIFETIME" default:"5m"`

	TraceExporter      string `env:"TRACE_EXPORTER" default:"otlp"` // "otlp" | "stdout" | ""
	TraceCollectorHost string `env:"TRACE_COLLECTOR_HOST"`
	TraceCollectorPort int    `env:"TRACE_COLLECTOR_PORT"`
}
