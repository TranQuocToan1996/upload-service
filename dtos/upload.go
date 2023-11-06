package dtos

import (
	"time"

	"upload_service/models"
)

type UploadFileResponse struct {
	Meta Meta `json:"meta"`
}

type DownloadFileRequest struct {
	FileID uint `param:"file_id" validator:"required"`
}

type DownloadFileErrResponse struct {
	Meta Meta `json:"meta"`
}

type GetListFilesUploadRequest struct {
	Limit               int64               `query:"limit" validate:"gte=0"`
	Offset              int64               `query:"offset" validate:"gte=0"`
	UserID              uint                `query:"user_id"`
	Status              models.UploadStatus `query:"status" validate:"omitempty,oneof=Pending Fail Done"`
	DestinationFileName string              `query:"destination_file_name"`
}

type GetListFilesUploadResponse struct {
	Meta Meta                             `json:"meta"`
	Data []GetListFilesUploadResponseData `json:"data,omitempty"`
}

type GetListFilesUploadResponseData struct {
	ID                  uint                `json:"id"`
	SourceFileName      string              `json:"source_file_name"`
	DestinationFileName string              `json:"destination_file_name"`
	ContentType         string              `json:"content_type"`
	SizeBytes           int64               `json:"size"`
	UserID              uint                `json:"user_id"`
	Status              models.UploadStatus `json:"status"`
	CreatedAt           time.Time           `json:"created_at"`
	UpdatedAt           time.Time           `json:"updated_at"`
}
