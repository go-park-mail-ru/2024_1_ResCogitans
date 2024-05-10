package avatar

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/go-park-mail-ru/2024_1_ResCogitans/internal/config"
	"github.com/go-park-mail-ru/2024_1_ResCogitans/internal/entities"
	"github.com/go-park-mail-ru/2024_1_ResCogitans/utils/logger"
)

var (
	maxFileSizeMB    = 8
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

	return fileType == "image/png" || fileType == "image/jpeg" || fileType == "image/jpg" || fileType == "image/webp"
}

func ValidateFileSize(handler *multipart.FileHeader) bool {
	// Get file size
	fileSize := handler.Size

	// Check if file size exceeds a certain limit
	return fileSize <= maxFileSizeBytes
}

func Upload(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path

	segments := strings.Split(path, "/")

	userID, err := strconv.Atoi(segments[len(segments)-2])
	if err != nil {
		http.Error(w, "Ошибка преобразования в число", http.StatusBadRequest)
		return
	}

	cfg, err := config.LoadConfig()
	if err != nil {
		http.Error(w, "Ошибка чтения файла конфигурации", http.StatusBadRequest)
		return
	}

	if r.Method != http.MethodPost {
		http.Error(w, "Метод не разрешен", http.StatusMethodNotAllowed)
		return
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Ошибка получения файла", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Validate
	flag := ValidateFileExtension(*header)
	if !flag {
		http.Error(w, "Неверный формат!", http.StatusBadRequest)
		return
	}

	flag = ValidateFileSize(header)
	if !flag {
		http.Error(w, "Слишком большой размер файла!", http.StatusBadRequest)
		return
	}

	// sh := sha256.New()
	// sh.Write([]byte(header.Filename))
	// hashBytes1 := sh.Sum(nil)
	// newPath := hex.EncodeToString(hashBytes1)

	uploadPath := filepath.Join(cfg.FileUploadPath, header.Filename)

	out, err := os.Create(uploadPath)
	if err != nil {
		http.Error(w, "Ошибка создания файла", http.StatusInternalServerError)
		return
	}
	defer out.Close()

	_, err = io.Copy(out, file)
	if err != nil {
		http.Error(w, "Ошибка сохранения файла", http.StatusInternalServerError)
		return
	}

	err = insertDataToDB(userID, header.Filename)
	if err != nil {
		http.Error(w, "Ошибка добавления пути в базу", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Файл успешно сохранен"))
}

func insertDataToDB(userID int, path string) error {
	log := logger.Logger()
	url := fmt.Sprintf("http://localhost:8080/api/profile/%d/edit", userID)

	var data entities.UserProfile
	data.Avatar = path

	jsonData, err := json.Marshal(data)
	if err != nil {
		log.Error("Ошибка при преобразовании в JSON:", err)
		return err
	}

	// Отправка POST-запроса
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		log.Error("Ошибка при отправке запроса:", err)
		return err
	}
	defer resp.Body.Close()

	// Чтение ответа
	_, err = io.ReadAll(resp.Body)
	if err != nil {
		log.Error("Ошибка при чтении ответа:", err)
		return err
	}
	return nil
}
