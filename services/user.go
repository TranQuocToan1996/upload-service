package services

import (
	"errors"

	"upload_service/dtos"
	"upload_service/models"
	"upload_service/repositories"
	tokens "upload_service/token"
)

var ErrUpdateErr = errors.New("update error")

type UserService interface {
	Create(newUser *models.User) error
	GetByUserName(userName string) (*models.User, error)
	GetByID(id uint) (*models.User, error)
	RevokeToken(request dtos.RevokeTokenRequest,
		claims *tokens.UserClaims) (*dtos.RevokeTokenResponse, error)
}

func ProvideUserService(
	userRepo repositories.UserRepository,
) UserService {
	return &userService{
		userRepo: userRepo,
	}
}

type userService struct {
	userRepo repositories.UserRepository
}

func (s *userService) Create(newUser *models.User) error {
	return s.userRepo.Create(newUser)
}

func (s *userService) GetByUserName(userName string) (*models.User, error) {
	return s.userRepo.GetByUserName(userName)
}

func (s *userService) GetByID(id uint) (*models.User, error) {
	return s.userRepo.GetByID(id)
}

func (s *userService) RevokeToken(request dtos.RevokeTokenRequest,
	claims *tokens.UserClaims,
) (*dtos.RevokeTokenResponse, error) {
	response := &dtos.RevokeTokenResponse{Meta: dtos.GetMeta(dtos.InternalError)}

	user, err := s.userRepo.GetByID(claims.UserID)
	if err != nil {
		return response, err
	}

	if request.RevokeTokenAt <= user.RevokeTokenAt {
		return response, err
	}

	rowsAffected, err := s.userRepo.UpdateRevokeAt(claims.UserID, request.RevokeTokenAt)
	if err != nil || rowsAffected <= 0 {
		return response, ErrUpdateErr
	}

	response.Meta = dtos.GetMeta(dtos.Success)
	return response, nil
}
