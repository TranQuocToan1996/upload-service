package routes

import (
	"upload_service/cmd/api/handlers"
	"upload_service/middlewares"
	tokens "upload_service/token"

	"github.com/golang-jwt/jwt/v5"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/ory/viper"
)

func Router(authHandler *handlers.AuthHandler) *echo.Echo {
	r := echo.New()
	r.Use(middleware.Logger())
	r.Use(middlewares.CORS())
	r.Use(middleware.Secure())
	r.Use(middleware.Recover())

	// Public routes
	publicGroup := r.Group("upload-service/v1")

	jwtConfig := echojwt.Config{
		NewClaimsFunc: func(c echo.Context) jwt.Claims {
			return new(tokens.JWTClaims)
		},
		KeyFunc: func(token *jwt.Token) (interface{}, error) {
			return []byte(viper.GetString("KEY")), nil
		},
	}

	privateGroup := publicGroup.Group("")
	privateGroup.Use(echojwt.WithConfig(jwtConfig))

	userGroup := publicGroup.Group("/users")
	userGroup.POST("", authHandler.Register)
	return r
}
