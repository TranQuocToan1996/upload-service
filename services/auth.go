package services

import (
	"encoding/hex"
	"errors"

	"upload_service/config"
	"upload_service/dtos"
	"upload_service/models"
	tokens "upload_service/token"
	"upload_service/utils"
)

var (
	ErrUserExist     = errors.New("user already created")
	ErrUserNotExist  = errors.New("user not exist")
	ErrPasswordWrong = errors.New("password wrong")
)

type AuthService interface {
	Register(request dtos.RegisterRequest) (*dtos.RegisterResponse, error)
	Login(request dtos.LoginRequest) (*dtos.LoginResponse, error)
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

func (s *authService) Register(
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

	token, err := s.createToken(&tokens.UserClaims{
		UserID: newUser.ID,
	})
	if err != nil {
		return response, err
	}

	response.Data = &dtos.Tokens{
		AccessToken: token,
	}
	response.Meta = dtos.GetMeta(dtos.Success)

	return response, nil
}

func (s *authService) Login(
	request dtos.LoginRequest,
) (*dtos.LoginResponse, error) {
	response := &dtos.LoginResponse{Meta: dtos.GetMeta(dtos.InternalError)}
	existUser, _ := s.userService.GetByUserName(request.UserName)
	if existUser == nil {
		response.Meta = dtos.GetMeta(dtos.UserNotExist)
		return response, ErrUserNotExist
	}

	salt, err := hex.DecodeString(existUser.Salt)
	if err != nil {
		return response, err
	}

	match := utils.IsPasswordsMatch(existUser.Password, request.Password, salt)
	if !match {
		response.Meta = dtos.GetMeta(dtos.PasswordWrong)
		return response, ErrPasswordWrong
	}

	token, err := s.createToken(&tokens.UserClaims{
		UserID: existUser.ID,
	})
	if err != nil {
		return response, err
	}

	response.Data = &dtos.Tokens{
		AccessToken: token,
	}
	response.Meta = dtos.GetMeta(dtos.Success)

	return response, nil
}

func (s *authService) createToken(claims *tokens.UserClaims) (string, error) {
	return s.tokenProvider.CreateToken(claims)
}
