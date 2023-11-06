package repositories

import (
	"upload_service/models"

	"gorm.io/gorm"
)

type UploadRepository interface {
	Create(upload *models.Upload) error
	UpdateStatus(uploadID uint, status models.UploadStatus) (int64, error)
	GetByID(id uint) (*models.Upload, error)
}

type uploadRepository struct {
	db *gorm.DB
}

func ProvideUploadRepository(db *gorm.DB) UploadRepository {
	return &uploadRepository{db: db}
}

func (r *uploadRepository) Create(upload *models.Upload) error {
	return r.db.Create(upload).Error
}

func (r *uploadRepository) UpdateStatus(uploadID uint, status models.UploadStatus) (int64, error) {
	res := r.db.Model(&models.Upload{}).
		Where("id = ?", uploadID).
		Update("status", status)
	return res.RowsAffected, res.Error
}

func (r *uploadRepository) GetByID(id uint) (*models.Upload, error) {
	var upload models.Upload
	err := r.db.First(&upload, id).Error
	if err != nil {
		return nil, err
	}
	return &upload, nil
}
