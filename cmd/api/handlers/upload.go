package handlers

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"upload_service/config"
	"upload_service/dtos"
	baseHanlder "upload_service/handlers"
	"upload_service/models"
	"upload_service/services"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"github.com/spf13/viper"
)

type UploadHandler struct {
	baseHanlder.BaseHandler

	validator     *validator.Validate
	config        config.Config
	uploadService services.UploadService
	userService   services.UserService
}

func ProvideUploadHandler(
	validator *validator.Validate,
	config config.Config,
	uploadService services.UploadService,
	userService services.UserService,
) *UploadHandler {
	return &UploadHandler{
		validator:     validator,
		config:        config,
		uploadService: uploadService,
		userService:   userService,
	}
}

func (h *UploadHandler) List(c echo.Context) error {
	var (
		request  = dtos.GetListFilesUploadRequest{}
		response = &dtos.GetListFilesUploadResponse{Meta: dtos.GetMeta(dtos.InternalError)}
	)

	err := c.Bind(&request)
	if err != nil {
		response.Meta = dtos.GetMeta(dtos.BindError)
		return c.JSON(h.GetHTTPCode(response.Meta.Code), response)
	}

	if err := h.validator.Struct(request); err != nil {
		response.Meta = dtos.GetMeta(dtos.BindError)
		return c.JSON(h.GetHTTPCode(response.Meta.Code), response)
	}

	claims, err := h.GetUserClaims(c)
	if err != nil {
		response.Meta = dtos.GetMeta(dtos.InternalError)
		return c.JSON(h.GetHTTPCode(response.Meta.Code), response)
	}

	if h.IsRevokeToken(h.userService, claims) {
		response.Meta = dtos.GetMeta(dtos.TokenRevoke)
		return c.JSON(h.GetHTTPCode(response.Meta.Code), response)
	}

	response, err = h.uploadService.ListFiles(request)
	if err != nil {
		log.Errorf("[ListImages] err: %v", err)
		return c.JSON(h.GetHTTPCode(response.Meta.Code), response)
	}

	return c.JSON(http.StatusOK, response)
}

func (h *UploadHandler) Download(c echo.Context) error {
	var (
		request  = dtos.DownloadFileRequest{}
		response = &dtos.DownloadFileErrResponse{Meta: dtos.GetMeta(dtos.InternalError)}
	)

	err := c.Bind(&request)
	if err != nil {
		response.Meta = dtos.GetMeta(dtos.BindError)
		return c.JSON(h.GetHTTPCode(response.Meta.Code), response)
	}

	if err := h.validator.Struct(request); err != nil {
		response.Meta = dtos.GetMeta(dtos.BindError)
		return c.JSON(h.GetHTTPCode(response.Meta.Code), response)
	}

	claims, err := h.GetUserClaims(c)
	if err != nil {
		response.Meta = dtos.GetMeta(dtos.InternalError)
		return c.JSON(h.GetHTTPCode(response.Meta.Code), response)
	}

	if h.IsRevokeToken(h.userService, claims) {
		response.Meta = dtos.GetMeta(dtos.TokenRevoke)
		return c.JSON(h.GetHTTPCode(response.Meta.Code), response)
	}

	data, contentType, err := h.uploadService.DownloadByID(request.FileID)
	if err != nil {
		log.Errorf("[DownloadImage] err: %v", err)
		return c.JSON(h.GetHTTPCode(response.Meta.Code), response)
	}

	return c.Blob(http.StatusOK, contentType, data)
}

func (h *UploadHandler) Upload(c echo.Context) error {
	var (
		response = &dtos.UploadFileResponse{Meta: dtos.GetMeta(dtos.InternalError)}
		err      error
	)
	defer func() {
		if err != nil {
			log.Errorf("[Upload] err: %v", err)
		}
	}()

	claims, err := h.GetUserClaims(c)
	if err != nil {
		return c.JSON(h.GetHTTPCode(response.Meta.Code), response)
	}

	if h.IsRevokeToken(h.userService, claims) {
		response.Meta = dtos.GetMeta(dtos.TokenRevoke)
		return c.JSON(h.GetHTTPCode(response.Meta.Code), response)
	}

	file, err := c.FormFile("data")
	if err != nil {
		return c.JSON(h.GetHTTPCode(response.Meta.Code), response)
	}

	upload := models.Upload{
		Status:              models.Pending,
		UserID:              claims.UserID,
		ContentType:         file.Header.Get("Content-Type"),
		SizeBytes:           file.Size,
		SourceFileName:      file.Filename,
		DestinationFileName: h.getDestinationName(file.Filename),
	}
	err = h.uploadService.Create(&upload)
	if err != nil {
		return c.JSON(h.GetHTTPCode(response.Meta.Code), response)
	}

	status := models.Fail
	defer func() {
		err := h.uploadService.UpdateStatus(upload.ID, status)
		if err != nil {
			log.Errorf("[Upload] err: %v", err)
		}
	}()

	if !h.isValidContentType(upload.ContentType) || !h.isValidSize(upload.SizeBytes) {
		return c.JSON(h.GetHTTPCode(response.Meta.Code), response)
	}

	source, err := file.Open()
	if err != nil {
		return c.JSON(h.GetHTTPCode(response.Meta.Code), response)
	}

	// ? May run in goroutine. And add another GET API uploadStatus for client polling
	err = h.write(source, upload.DestinationFileName)
	if err != nil {
		return c.JSON(h.GetHTTPCode(response.Meta.Code), response)
	}

	// update status done in above defer func
	status = models.Done

	response.Meta = dtos.GetMeta(dtos.Success)
	return c.JSON(http.StatusOK, response)
}

func (h *UploadHandler) write(source io.ReadCloser, newName string) error {
	defer source.Close()
	if _, err := os.Stat(h.config.PathUpload); os.IsNotExist(err) {
		os.MkdirAll(h.config.PathUpload, os.ModePerm)
	}
	destination, err := os.Create(fmt.Sprintf("%v/%v", h.config.PathUpload, newName))
	if err != nil {
		return err
	}
	defer destination.Close()
	if _, err = io.Copy(destination, source); err != nil {
		return err
	}
	return nil
}

func (h *UploadHandler) isValidContentType(contentType string) bool {
	listRaw := viper.GetString("WHITE_LIST_CONTENT_TYPE")
	if len(listRaw) == 0 {
		return true
	}
	lists := strings.Split(listRaw, ",")
	for _, compare := range lists {
		if strings.EqualFold(contentType, compare) {
			return true
		}
	}

	return false
}

// Assume same file name can upload multiple times
// so that add UNIX time upload to track
// and avoid truncate when os.Create(fileName) in /tmp
func (h *UploadHandler) getDestinationName(fileName string) string {
	return fmt.Sprintf("%v_%v", time.Now().Unix(), fileName)
}

func (h *UploadHandler) isValidSize(curSize int64) bool {
	return curSize <= h.config.UploadLimitBytes
}
