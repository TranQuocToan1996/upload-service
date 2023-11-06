package dtos

type UploadFileResponse struct {
	Meta Meta `json:"meta"`
}

type DownloadFileRequest struct {
	FileID uint `param:"file_id" validator:"required"`
}

type DownloadFileErrResponse struct {
	Meta Meta `json:"meta"`
}
