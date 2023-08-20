package env

// Env はアプリケーションで使用する環境変数
type Env struct {
	API   API
	MySQL MySQL
}

// API は API に関する環境変数
type API struct {
	Host string `envconfig:"API_HOST" required:"true"`
	Port int    `envconfig:"API_PORT" required:"true"`
}

// MySQL は MySQL に関する環境変数
type MySQL struct {
	Host     string `envconfig:"MYSQL_HOST" required:"true"`
	Port     int    `envconfig:"MYSQL_PORT" required:"true"`
	Database string `envconfig:"MYSQL_DATABASE" required:"true"`
	User     string `envconfig:"MYSQL_USER" required:"true"`
	Password string `envconfig:"MYSQL_PASSWORD" required:"true"`
}
