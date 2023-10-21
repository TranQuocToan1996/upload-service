package models

import "upload_service/repositories/db"

type UploadStatus string

const (
	Pending UploadStatus = "Pending"
	Fail    UploadStatus = "Fail"
	Done    UploadStatus = "Done"
)

type Upload struct {
	db.BaseModel

	SourceFileName      string       `gorm:"column:source_file_name"`
	DestinationFileName string       `gorm:"column:destination_file_name"`
	ContentType         string       `gorm:"column:content_type"`
	SizeBytes           int64        `gorm:"column:size"`
	UserID              uint         `gorm:"column:user_id"`
	Status              UploadStatus `gorm:"column:status"`
}

func (Upload) TableName() string {
	return "uploads"
}
