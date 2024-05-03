package avatar

import (
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-park-mail-ru/2024_1_ResCogitans/internal/config"
)

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

	// Validate
	flag := ValidateFileExtension(*header)
	if !flag {
		http.Error(w, "Not valid format!", http.StatusBadRequest)
		return
	}

	flag = ValidateFileSize(header)
	if !flag {
		http.Error(w, "Not valid format!", http.StatusBadRequest)
		return
	}

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
