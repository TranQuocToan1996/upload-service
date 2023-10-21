package handlers

import (
	"errors"
	"net/http"
	"strings"

	"upload_service/services"
	tokens "upload_service/token"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"github.com/spf13/viper"
)

var ErrEmptyClaims = errors.New("empty claims")

type BaseHandler struct{}

func (b *BaseHandler) GetHTTPCode(code string) int {
	const (
		lengthHTTPCode = 3
	)
	if len(code) < lengthHTTPCode {
		return http.StatusInternalServerError
	}

	switch {
	case strings.HasPrefix(code, "200"):
		return http.StatusOK
	case strings.HasPrefix(code, "400"):
		return http.StatusBadRequest
	case strings.HasPrefix(code, "403"):
		return http.StatusForbidden
	case strings.HasPrefix(code, "404"):
		return http.StatusNotFound
	case strings.HasPrefix(code, "500"):
		return http.StatusInternalServerError
	default:
		return http.StatusInternalServerError
	}
}

func (b *BaseHandler) getClaims(c echo.Context) (jwt.Claims, error) {
	defaultKey := viper.GetString("KEY_CLAIMS")
	if len(defaultKey) == 0 {
		defaultKey = "user" // Default setting of library
	}

	if c.Get(defaultKey) == nil {
		return nil, ErrEmptyClaims
	}

	tokenObj, ok := c.Get(defaultKey).(*jwt.Token)
	if !ok {
		return nil, ErrEmptyClaims
	}

	return tokenObj.Claims, nil
}

func (b *BaseHandler) GetUserClaims(c echo.Context) (*tokens.JWTClaims, error) {
	claimsObj, err := b.getClaims(c)
	if err != nil {
		return nil, err
	}

	claims, ok := claimsObj.(*tokens.JWTClaims)
	if !ok {
		return nil, ErrEmptyClaims
	}

	return claims, nil
}

func (b *BaseHandler) IsRevokeToken(userService services.UserService, claims *tokens.JWTClaims) bool {
	user, _ := userService.GetByID(claims.UserID)
	if user == nil {
		return true
	}
	if user.RevokeTokenAt == 0 {
		return false
	}
	return user.RevokeTokenAt >= claims.IssuedAt.Unix()
}
