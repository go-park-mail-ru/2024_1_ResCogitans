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
	magicTable       = map[string]string{
		"\xff\xd8\xff":      "image/jpeg",
		"\x89PNG\r\n\x1a\n": "image/png",
		"GIF87a":            "image/gif",
		"GIF89a":            "image/gif",
	}
)

func DetectType(b []byte) bool {
	flag := false
	s := string(b)
	for key, val := range magicTable {
		if strings.HasPrefix(s, key) {
			fmt.Println(val)
			flag = true
		}
	}
	return flag
}

func ValidateFileExtension(file multipart.File) bool {
	buff := make([]byte, 512)
	if _, err := file.Read(buff); err != nil {
		return false
	}

	val := DetectType(buff)

	return val
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
	defer file.Close()

	// Validate
	flag := ValidateFileExtension(file)
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
