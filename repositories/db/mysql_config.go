package db

import "fmt"

// MySQLConfig define MySQLConfig structure
type MySQLConfig struct {
	Username    string
	Password    string
	Host        string
	Port        int64
	Database    string
	MaxOpen     int
	MaxIdle     int
	MaxLifetime int
	EnableDebug bool
}

// DriverName return driver name
func (t *MySQLConfig) DriverName() string {
	return "mysql"
}

// ConnectionString return connection string
func (t *MySQLConfig) ConnectionString() string {
	return fmt.Sprintf(
		"%v:%v@tcp(%v:%v)/%v?parseTime=true&charset=utf8mb4",
		t.Username,
		t.Password,
		t.Host,
		t.Port,
		t.Database,
	)
}

// MaxOpenConns config limit pool connection
func (t *MySQLConfig) MaxOpenConns() int {
	return t.MaxOpen
}

// MaxIdleConns config max idle connection
func (t *MySQLConfig) MaxIdleConns() int {
	return t.MaxIdle
}

// MaxLifeTime config max left time connection
func (t *MySQLConfig) MaxLifeTime() int {
	return t.MaxLifetime
}

// IsDebugMode on/off debug ORM
func (t *MySQLConfig) IsDebugMode() bool {
	return t.EnableDebug
}
