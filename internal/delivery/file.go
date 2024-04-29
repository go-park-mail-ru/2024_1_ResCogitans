package delivery

import (
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"strings"

	"github.com/go-park-mail-ru/2024_1_ResCogitans/internal/config"
	"github.com/go-park-mail-ru/2024_1_ResCogitans/utils/logger"
	"github.com/pkg/errors"
)

type FileHandler struct{}

type FileResponse struct {
	Path string `json:"path"`
}

var (
	maxFileSizeMB    = 1
	maxFileSizeBytes = int64(maxFileSizeMB) * 1024 * 1024
)

func ValidateFileExtension(fileHeader multipart.FileHeader) bool {
	buff := make([]byte, 512)
	file, err := fileHeader.Open()
	if err != nil {
		return false
	}
	defer file.Close()
	if _, err := file.Read(buff); err != nil {
		return false
	}

	fileType := strings.ToLower(http.DetectContentType(buff))

	return fileType == "image/png" || fileType == "image/jpeg" || fileType == "image/jpg"
}

func ValidateFileSize(handler *multipart.FileHeader) bool {
	// Get file size
	fileSize := handler.Size

	// Check if file size exceeds a certain limit
	return fileSize <= maxFileSizeBytes
}

func SaveFile(r *http.Request) (string, error) {
	logger := logger.Logger()

	r.ParseMultipartForm(10 << 20)
	file, handler, err := r.FormFile("file")
	if err != nil {
		logger.Error("Error while retrieving file:", "error", err)
		return string(""), err
	}
	fmt.Println("No")
	defer file.Close()

	// Validate
	flag := ValidateFileExtension(*handler)
	if !flag {
		logger.Error("Not valid format!")
		return string(""), errors.New("Not valid format!")
	}

	flag = ValidateFileSize(handler)
	if !flag {
		logger.Error("Error while checking file size:", "error", err)
		return string(""), errors.New("Error while checking file size!")
	}

	cfg, _ := config.LoadConfig()
	targetFile, err := os.Create(cfg.FileUploadPath + handler.Filename)
	if err != nil {
		logger.Error("Error while creating file:", "error", err)
		return string(""), err
	}
	defer targetFile.Close()

	_, err = io.Copy(targetFile, file)
	if err != nil {
		logger.Error("Error while creating file:", "error", err)
		return string(""), err
	}

	return targetFile.Name(), nil
}
