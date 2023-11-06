package services

import (
	"errors"
	"fmt"
	"os"

	"upload_service/config"
	"upload_service/dtos"
	"upload_service/models"
	"upload_service/repositories"

	"github.com/jinzhu/copier"
)

var ErrUploadFileNotDone = errors.New("upload file not done")

type UploadService interface {
	Create(upload *models.Upload) error
	UpdateStatus(uploadID uint, status models.UploadStatus) error
	DownloadByID(uploadID uint) ([]byte, string, error)
	ListFiles(request dtos.GetListFilesUploadRequest) (*dtos.GetListFilesUploadResponse, error)
}

func ProvideUploadService(
	config config.Config,
	uploadRepo repositories.UploadRepository,
) UploadService {
	return &uploadService{
		config:     config,
		uploadRepo: uploadRepo,
	}
}

type uploadService struct {
	config     config.Config
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

func (s *uploadService) DownloadByID(uploadID uint) ([]byte, string, error) {
	// ? assume a user can download image of another user
	upload, err := s.uploadRepo.GetByID(uploadID)
	if err != nil {
		return nil, "", err
	}
	// ? For pending upload, need cronjob check and update status before user can download
	if upload.Status != models.Done {
		return nil, "", ErrUploadFileNotDone
	}
	data, err := os.ReadFile(fmt.Sprintf("%v/%v", s.config.PathUpload, upload.DestinationFileName))
	if err != nil {
		return nil, "", err
	}
	return data, upload.ContentType, nil
}

func (s *uploadService) ListFiles(request dtos.GetListFilesUploadRequest) (*dtos.GetListFilesUploadResponse, error) {
	response := &dtos.GetListFilesUploadResponse{Meta: dtos.GetMeta(dtos.InternalError)}

	if request.Limit == 0 || request.Limit > s.config.MaxGetLimit {
		request.Limit = s.config.DefaultGetLimit
	}

	list, total, err := s.uploadRepo.GetByFilter(repositories.GetUploadFilter{
		Limit:               request.Limit,
		Offset:              request.Offset,
		UserID:              request.UserID,
		Status:              request.Status,
		DestinationFileName: request.DestinationFileName,
	})
	if err != nil {
		return nil, err
	}

	responseData := make([]dtos.GetListFilesUploadResponseData, 0, len(list))
	err = copier.Copy(&responseData, list)
	if err != nil {
		return nil, err
	}

	response.Data = responseData
	response.Meta = dtos.GetMeta(dtos.Success)
	response.Meta.SetLimit(request.Limit).SetOffset(request.Offset).SetTotal(total)

	return response, nil
}
