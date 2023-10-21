package services

import (
	"encoding/hex"
	"errors"

	"upload_service/config"
	"upload_service/dtos"
	"upload_service/models"
	tokens "upload_service/token"
	"upload_service/utils"

	"github.com/labstack/echo/v4"
)

var ErrUserExist = errors.New("user already created")

type AuthService interface {
	Register(c echo.Context, request dtos.RegisterRequest) (*dtos.RegisterResponse, error)
}

func ProvideAuthService(
	config config.Config,
	userService UserService,
	tokenProvider tokens.TokenProvider,
) AuthService {
	return &authService{
		tokenProvider: tokenProvider,
		userService:   userService,
		config:        config,
	}
}

type authService struct {
	config        config.Config
	tokenProvider tokens.TokenProvider
	userService   UserService
}

func (s *authService) Register(c echo.Context,
	request dtos.RegisterRequest,
) (*dtos.RegisterResponse, error) {
	response := &dtos.RegisterResponse{Meta: dtos.GetMeta(dtos.InternalError)}
	existUser, _ := s.userService.GetByUserName(request.UserName)
	if existUser != nil {
		response.Meta = dtos.GetMeta(dtos.UserExist)
		return response, ErrUserExist
	}

	salt, err := utils.GenerateRandomSalt(s.config.SaltLength)
	if err != nil {
		return response, err
	}

	hashPassword := utils.HashPassword(request.Password, salt)

	newUser := models.User{
		UserName: request.UserName,
		Password: hashPassword,
		Salt:     hex.EncodeToString(salt),
	}

	err = s.userService.Create(&newUser)
	if err != nil {
		return response, err
	}

	token := s.createToken(&tokens.UserClaims{
		UserID: newUser.ID,
	})

	response.Data = &dtos.Tokens{
		AccessToken: token,
	}
	response.Meta = dtos.GetMeta(dtos.Success)

	return response, nil
}

func (s *authService) createToken(claims *tokens.UserClaims) string {
	return s.tokenProvider.CreateToken(claims)
}
