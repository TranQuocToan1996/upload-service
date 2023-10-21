package handlers

import (
	"net/http"

	"upload_service/dtos"
	baseHanlder "upload_service/handlers"
	"upload_service/services"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
)

type AuthHandler struct {
	baseHanlder.BaseHandler

	validator   *validator.Validate
	authService services.AuthService
}

func ProvideAuthHandler(
	validator *validator.Validate,
	authService services.AuthService,
) *AuthHandler {
	return &AuthHandler{
		validator:   validator,
		authService: authService,
	}
}

func (h *AuthHandler) Register(c echo.Context) error {
	var (
		request  = dtos.RegisterRequest{}
		response = &dtos.RegisterResponse{}
	)

	err := c.Bind(&request)
	if err != nil {
		response.Meta = dtos.GetMeta(dtos.BindError)
		return c.JSON(h.GetHTTPCode(response.Meta.Code), response)
	}

	if err := h.validator.Struct(request); err != nil {
		response.Meta = dtos.GetMeta(dtos.BindError)
		return c.JSON(h.GetHTTPCode(response.Meta.Code), response)
	}

	response, err = h.authService.Register(c, request)
	if err != nil {
		log.Errorf("[Register] err: %v", err)
		return c.JSON(h.GetHTTPCode(response.Meta.Code), response)
	}

	return c.JSON(http.StatusOK, response)
}
