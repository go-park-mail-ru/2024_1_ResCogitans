package usecase

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-park-mail-ru/2024_1_ResCogitans/internal/config"
)

type DownloadResponse struct {
	Href string `json:"href"`
}

func GetDownloadLink(filePath string) (string, error) {
	path := "jantugan/album/" + filePath
	url := "https://cloud-api.yandex.net/v1/disk/resources/download?path=" + path

	request, err := http.NewRequest("GET", url, nil)
	if err != nil {

		return "", err
	}
	cfg, _ := config.LoadConfig()

	request.Header.Set("Authorization", "OAuth "+cfg.Drive.Token)

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return "", err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return "", fmt.Errorf("download link request failed with status: %d", response.StatusCode)
	}

	var downloadResponse DownloadResponse
	if err := json.NewDecoder(response.Body).Decode(&downloadResponse); err != nil {
		return "", err
	}

	return downloadResponse.Href, nil
}
