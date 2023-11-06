package repositories

import (
	"upload_service/models"

	"gorm.io/gorm"
)

type UploadRepository interface {
	Create(upload *models.Upload) error
	UpdateStatus(uploadID uint, status models.UploadStatus) (int64, error)
	GetByID(id uint) (*models.Upload, error)
	GetByFilter(filter GetUploadFilter) ([]models.Upload, int64, error)
}

type GetUploadFilter struct {
	Limit               int64
	Offset              int64
	UserID              uint
	Status              models.UploadStatus
	DestinationFileName string
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

func (r *uploadRepository) GetByFilter(filter GetUploadFilter) ([]models.Upload, int64, error) {
	var (
		total   int64 = 0
		uploads []models.Upload
	)

	db := r.db.Model(&models.Upload{}).Order("id DESC")

	if len(filter.Status) > 0 {
		db = db.Where("status = ?", filter.Status)
	}
	if filter.UserID > 0 {
		db = db.Where("id = ?", filter.UserID)
	}
	if len(filter.DestinationFileName) > 0 {
		db = db.Where("destination_file_name LIKE ?", "%"+filter.DestinationFileName+"%")
	}

	db = db.Limit(int(filter.Limit)).Offset(int(filter.Offset))

	err := db.Find(&uploads).Error
	if err != nil {
		return nil, 0, err
	}

	err = db.Limit(-1).Offset(-1).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	return uploads, total, nil
}
