package db

import (
	"errors"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

// DB define db client
type DB = gorm.DB

// BaseModel defines base model
type BaseModel struct {
	ID        uint           `gorm:"column:id;primary_key;AUTO_INCREMENT"`
	CreatedAt time.Time      `gorm:"column:created_at;"`
	UpdatedAt time.Time      `gorm:"column:updated_at;"`
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at;"`
}

// ErrRecordNotFound ..
var ErrRecordNotFound = gorm.ErrRecordNotFound

// Open connection
func Open(config Config) (*DB, error) {
	logLevel := logger.Error
	if config.IsDebugMode() {
		logLevel = logger.Info
	}

	db, err := gorm.Open(mysql.Open(config.ConnectionString()), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
		Logger: logger.Default.LogMode(logLevel),
	})
	if err != nil {
		return nil, err
	}

	mysqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	mysqlDB.SetMaxOpenConns(config.MaxOpenConns())
	mysqlDB.SetMaxIdleConns(config.MaxIdleConns())
	mysqlDB.SetConnMaxLifetime(time.Duration(config.MaxLifeTime()) * time.Second)
	return db, nil
}

func IsRecordNotFoundErr(err error) bool {
	return errors.Is(err, ErrRecordNotFound)
}
