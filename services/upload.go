package services

import (
	"upload_service/models"
	"upload_service/repositories"
)

type UploadService interface {
	Create(upload *models.Upload) error
	UpdateStatus(uploadID uint, status models.UploadStatus) error
}

func ProvideUploadService(
	uploadRepo repositories.UploadRepository,
) UploadService {
	return &uploadService{
		uploadRepo: uploadRepo,
	}
}

type uploadService struct {
	uploadRepo repositories.UploadRepository
}

func (s *uploadService) Create(upload *models.Upload) error {
	return s.uploadRepo.Create(upload)
}

func (s *uploadService) UpdateStatus(uploadID uint, status models.UploadStatus) error {
	rowsAffected, err := s.uploadRepo.UpdateStatus(uploadID, status)
	if err != nil || rowsAffected <= 0 {
		return ErrUpdateErr
	}
	return nil
}
