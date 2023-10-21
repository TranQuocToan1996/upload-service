package main

import (
	"context"
	"net"
	"net/http"
	"time"

	"upload_service/cmd/api/handlers"
	"upload_service/cmd/api/routes"
	"upload_service/config"
	"upload_service/repositories"
	"upload_service/repositories/db"
	"upload_service/services"
	tokens "upload_service/token"
	"upload_service/utils"

	log "github.com/sirupsen/logrus"

	"github.com/go-playground/validator/v10"
	_ "github.com/go-sql-driver/mysql"
	"github.com/ory/viper"
	"go.uber.org/fx"
	"gorm.io/gorm"
)

func NewHTTPServer(
	lc fx.Lifecycle,
	authHanlder *handlers.AuthHandler,
) *http.Server {
	handler := routes.Router(authHanlder).Server.Handler
	server := &http.Server{
		Addr: viper.GetString("APP_HTTP_SERVER"), Handler: handler,
		ReadHeaderTimeout: time.Second * time.Duration(viper.GetInt("HTTP_READ_HEADER_TIMEOUT")),
	}
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			ln, err := net.Listen("tcp", server.Addr)
			if err != nil {
				return err
			}
			log.Println("Starting HTTP server at", server.Addr)
			go func() {
				err := server.Serve(ln)
				if err != nil {
					log.Panic(err)
				}
			}()

			return nil
		},
		OnStop: func(ctx context.Context) error {
			return server.Shutdown(ctx)
		},
	})

	return server
}

func main() {
	utils.InitConfig(".env")

	fx.New(
		fx.Provide(
			ProvideValidator,
			ProviveGormMySQL,
			NewHTTPServer,
			tokens.NewJWTProvider,
			config.ProvideConfig,
			repositories.ProvideUserRepository,
			services.ProvideAuthService,
			services.ProvideUserService,
			handlers.ProvideAuthHandler,
		),
		fx.Invoke(func(*http.Server) {}),
	).Run()
}

func ProviveGormMySQL() *gorm.DB {
	gormDB, err := db.Open(MySQLConfig())
	if err != nil {
		log.Fatalf("Connecting to MySQL: %v", err)
	}
	return gormDB
}

func ProvideValidator() *validator.Validate {
	return validator.New()
}

func MySQLConfig() db.Config {
	return &db.MySQLConfig{
		Username:    viper.GetString("MYSQL_USERNAME"),
		Password:    viper.GetString("MYSQL_PASSWORD"),
		Host:        viper.GetString("MYSQL_HOST"),
		Port:        viper.GetInt64("MYSQL_PORT"),
		Database:    viper.GetString("MYSQL_DATABASE"),
		MaxOpen:     viper.GetInt("MYSQL_POOL_SIZE"),
		MaxIdle:     viper.GetInt("MYSQL_MAX_IDLE"),
		MaxLifetime: viper.GetInt("MYSQL_MAX_LEFTIME"),
		EnableDebug: viper.GetBool("MYSQL_DEBUG"),
	}
}
