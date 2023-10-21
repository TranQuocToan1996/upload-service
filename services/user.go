package services

import (
	"upload_service/models"
	"upload_service/repositories"
)

type UserService interface {
	Create(newUser *models.User) error
	GetByUserName(userName string) (*models.User, error)
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
