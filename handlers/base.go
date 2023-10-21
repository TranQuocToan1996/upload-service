package handlers

import (
	"errors"
	"net/http"
	"strings"
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
