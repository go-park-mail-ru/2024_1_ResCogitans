package delivery

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

func Upload(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	file, header, err := r.FormFile("image")
	if err != nil {
		http.Error(w, "Error retrieving the file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	cwd, err := os.Getwd()
	if err != nil {
		http.Error(w, "Error getting current directory", http.StatusInternalServerError)
		return
	}

	uploadPath := filepath.Join(cwd, "../../../Modules/public", header.Filename)

	out, err := os.Create(uploadPath)
	if err != nil {
		fmt.Println(err)
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
