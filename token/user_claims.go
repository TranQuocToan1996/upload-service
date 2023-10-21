package tokens

import (
	"github.com/golang-jwt/jwt/v5"
)

type JWTClaims struct {
	jwt.RegisteredClaims
	UserClaims
}

type UserClaims struct {
	UserID uint `json:"user_id"`
}
