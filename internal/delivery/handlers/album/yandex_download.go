package album

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/go-park-mail-ru/2024_1_ResCogitans/internal/config"
	"github.com/go-park-mail-ru/2024_1_ResCogitans/utils/logger"
	"github.com/pkg/errors"
)

type UploadResponse struct {
	Href string `json:"href"`
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

	return fileType == "image/png" || fileType == "image/jpeg" || fileType == "image/jpg" || fileType == "image/webp"
}

func ValidateFileSize(handler *multipart.FileHeader) bool {
	// Get file size
	fileSize := handler.Size

	// Check if file size exceeds a certain limit
	return fileSize <= maxFileSizeBytes
}

func getURL(path, token string) (string, error) {
	uploadPath := "jantugan/album/" + path
	url := "https://cloud-api.yandex.net/v1/disk/resources/upload?path=" + uploadPath

	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", err
	}
	request.Header.Set("Authorization", "OAuth "+token)

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return "", err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return "", fmt.Errorf("getting URL failed with status: %d", response.StatusCode)
	}

	var uploadResponse UploadResponse
	if err := json.NewDecoder(response.Body).Decode(&uploadResponse); err != nil {
		return "", err
	}

	return uploadResponse.Href, nil
}

func uploadFile(file multipart.File, handler *multipart.FileHeader) (string, error) {
	logger := logger.Logger()

	// Валидация
	flag := ValidateFileExtension(*handler)
	if !flag {
		logger.Error("Not valid format!")
		return "", errors.New("Not valid format!")
	}

	flag = ValidateFileSize(handler)
	if !flag {
		logger.Error("Error while checking file size:", "error")
		return "", errors.New("Error while checking file size!")
	}

	// Берем хэш из названия файла
	sh := sha256.New()
	sh.Write([]byte(handler.Filename))
	hashBytes1 := sh.Sum(nil)
	newPath := hex.EncodeToString(hashBytes1)

	cfg, _ := config.LoadConfig()

	uploadURL, err := getURL(newPath, cfg.Drive.Token)
	if err != nil {
		return "", err
	}

	// Создаем HTTP запрос для загрузки файла
	request, err := http.NewRequest("PUT", uploadURL, file)
	if err != nil {
		return "", err
	}
	request.ContentLength = handler.Size

	// Отправляем запрос
	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return "", err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusCreated {
		logger.Error("upload failed with status: %d", response.StatusCode)
	}

	logger.Info("File uploaded successfully!")

	return newPath, nil
}

func insertDataToDB(albumID int, path string) error {
	log := logger.Logger()
	url := fmt.Sprintf("http://localhost:8080/album/%d/add", albumID) // Предполагаем, что ваше приложение слушает на порту 8080

	jsonData, err := json.Marshal(path)
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
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Error("Ошибка при чтении ответа:", err)
		return err
	}

	// Вывод ответа
	fmt.Println("Ответ от сервера:", string(body))
	return nil
}

func UploadImageAndInsert(w http.ResponseWriter, r *http.Request) {
	logger := logger.Logger()

	r.ParseMultipartForm(10 << 20)
	file, handler, err := r.FormFile("file")
	if err != nil {
		logger.Error("Error while retrieving file:", "error", err)
		http.Error(w, "Error while retrieving file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	albumID, err := strconv.Atoi(chi.URLParam(r, "albumID"))
	if err != nil {
		logger.Error("Cannot convert to int", err)
		http.Error(w, "Cannot convert to int", http.StatusBadRequest)
		return
	}

	path, err := uploadFile(file, handler)
	if err != nil {
		logger.Error("Error while uploading file:", "error", err)
		http.Error(w, "Error while uploading file", http.StatusBadRequest)
		return
	}

	err = insertDataToDB(albumID, path)
	if err != nil {
		logger.Error("Error updating DB", "error", err)
		http.Error(w, "Error updating DB", http.StatusBadRequest)
		return
	}
}
