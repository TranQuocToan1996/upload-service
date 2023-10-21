package db

// Config define config interface
type Config interface {
	DriverName() string
	ConnectionString() string
	IsDebugMode() bool
	MaxOpenConns() int
	MaxIdleConns() int
	MaxLifeTime() int
}
