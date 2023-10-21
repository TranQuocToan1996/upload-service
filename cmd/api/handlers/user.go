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

type UserHandler struct {
	baseHanlder.BaseHandler

	validator   *validator.Validate
	userService services.UserService
}

func ProvideUserHandler(
	validator *validator.Validate,
	userService services.UserService,
) *UserHandler {
	return &UserHandler{
		validator:   validator,
		userService: userService,
	}
}

func (h *UserHandler) RevokeToken(c echo.Context) error {
	var (
		request  = dtos.RevokeTokenRequest{}
		response = &dtos.RevokeTokenResponse{}
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

	claims, err := h.GetUserClaims(c)
	if err != nil {
		response.Meta = dtos.GetMeta(dtos.InternalError)
		return c.JSON(h.GetHTTPCode(response.Meta.Code), response)
	}

	// ? Move to middlewares
	if h.IsRevokeToken(h.userService, claims) {
		response.Meta = dtos.GetMeta(dtos.TokenRevoke)
		return c.JSON(h.GetHTTPCode(response.Meta.Code), response)
	}

	response, err = h.userService.RevokeToken(request, &claims.UserClaims)
	if err != nil {
		log.Errorf("[RevokeToken] err: %v", err)
		return c.JSON(h.GetHTTPCode(response.Meta.Code), response)
	}

	return c.JSON(http.StatusOK, response)
}
