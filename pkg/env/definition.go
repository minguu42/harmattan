package env

import "fmt"

// Env はアプリケーションで使用する環境変数
type Env struct {
	API   *API
	MySQL *MySQL
}

// API は API に関する環境変数
type API struct {
	Host string
	Port int
}

// MySQL は MySQL に関する環境変数
type MySQL struct {
	Host     string
	Port     int
	Database string
	User     string
	Password string
}

// DSN はデータベースとの接続に使用する Data Source Name を生成する
func (r *MySQL) DSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=true",
		r.User,
		r.Password,
		r.Host,
		r.Port,
		r.Database,
	)
}
