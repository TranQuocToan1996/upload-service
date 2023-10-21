package main

import (
	"database/sql"
	"embed"
	"errors"
	"fmt"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/pressly/goose/v3"
	"github.com/spf13/viper"
)

//go:embed sqls/*.sql
var embedMigrations embed.FS

func main() {
	viper.AutomaticEnv()
	viper.SetConfigFile(".env")
	err := viper.ReadInConfig()
	var pathErr *os.PathError
	if errors.As(err, &pathErr) {
		log.Printf("no config file '%s' not found. Using default values", ".env")
	} else if err != nil { // Handle other errors that occurred while reading the config file
		panic(fmt.Errorf("fatal error while reading the config file: %w", err))
	}
	var db *sql.DB
	// setup database
	db, err = goose.OpenDBWithDriver("mysql", fmt.Sprintf(
		"%v:%v@tcp(%v:%v)/%v?parseTime=true&charset=utf8mb4",
		viper.GetString("MYSQL_USERNAME"),
		viper.GetString("MYSQL_PASSWORD"),
		viper.GetString("MYSQL_HOST"),
		viper.GetInt64("MYSQL_PORT"),
		viper.GetString("MYSQL_DATABASE"),
	))
	if err != nil {
		log.Fatalf("goose: failed to open DB: %v\n", err)
	}

	defer func() {
		if err := db.Close(); err != nil {
			log.Fatalf("goose: failed to close DB: %v\n", err)
		}
	}()

	goose.SetBaseFS(embedMigrations)

	if err := goose.Up(db, "sqls"); err != nil {
		panic(err)
	}
}
