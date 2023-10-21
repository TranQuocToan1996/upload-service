package repositories

import (
	"upload_service/models"

	"gorm.io/gorm"
)

type UserRepository interface {
	Create(user *models.User) error
	GetByUserName(userName string) (*models.User, error)
	UpdateRevokeAt(userID uint, unixSeconds int64) (int64, error)
	GetByID(id uint) (*models.User, error)
}

type userRepository struct {
	db *gorm.DB
}

func ProvideUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(user *models.User) error {
	return r.db.Create(user).Error
}

func (r *userRepository) GetByUserName(userName string) (*models.User, error) {
	var user models.User
	err := r.db.First(&user, "user_name = ?", userName).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) UpdateRevokeAt(userID uint, unixSeconds int64) (int64, error) {
	res := r.db.Model(&models.User{}).
		Where("id = ?", userID).
		Update("revoke_token_at", unixSeconds)
	return res.RowsAffected, res.Error
}

func (r *userRepository) GetByID(id uint) (*models.User, error) {
	var user models.User
	err := r.db.First(&user, id).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}
