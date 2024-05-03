package avatar

import (
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/go-park-mail-ru/2024_1_ResCogitans/internal/config"
)

// Upload
// @Summary Загрузить файл
// @Description Загружает файл в указанный путь в конфигурации
// @Tags Фотографии
// @Accept  multipart/form-data
// @Produce json
// @Param   file formData file true  "Файл для загрузки"
// @Success 200 {string} string "Файл успешно загружен"
// @Failure 400 {object} httperrors.HttpError
// @Failure 405 {object} httperrors.HttpError
// @Failure 500 {object} httperrors.HttpError
// @Router /api/profile/{id}/upload [post]
func Upload(w http.ResponseWriter, r *http.Request) {
	cfg, err := config.LoadConfig()
	if err != nil {
		http.Error(w, "Failed to load config", http.StatusBadRequest)
		return
	}

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Error retrieving the file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	uploadPath := filepath.Join(cfg.FileUploadPath, header.Filename)

	out, err := os.Create(uploadPath)
	if err != nil {
		http.Error(w, "Error creating the file", http.StatusInternalServerError)
		return
	}
	defer out.Close()

	_, err = io.Copy(out, file)
	if err != nil {
		http.Error(w, "Error saving the file", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("File uploaded successfully"))
}
