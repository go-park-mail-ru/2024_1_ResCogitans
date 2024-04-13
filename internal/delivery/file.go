package delivery

import (
	"io"
	"net/http"
	"os"

	"github.com/go-park-mail-ru/2024_1_ResCogitans/internal/config"
	"github.com/go-park-mail-ru/2024_1_ResCogitans/utils/errors"
	"github.com/go-park-mail-ru/2024_1_ResCogitans/utils/logger"
)

type FileHandler struct{}

type FileResponse struct {
	Path string `json:"path"`
}

var (
	errUploadFile = errors.HttpError{
		Code:    http.StatusInternalServerError,
		Message: "failed upload file",
	}
)

func SaveFile(r *http.Request) (string, error) {
	logger := logger.Logger()

	r.ParseMultipartForm(10 << 20)
	file, handler, err := r.FormFile("file")
	if err != nil {
		logger.Error("Error while retrieving file:", "error", err)
		return string(""), err
	}
	defer file.Close()

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
