package usecase

import (
	"fmt"
	"net/http"

	"github.com/go-park-mail-ru/2024_1_ResCogitans/internal/config"
)

func deleteResource(path string) error {
	url := fmt.Sprintf("https://cloud-api.yandex.net/v1/disk/resources?path=%s&permanently=true", path)

	request, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return err
	}

	cfg, _ := config.LoadConfig()
	request.Header.Set("Authorization", "OAuth "+cfg.Drive.Token)

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("delete resource request failed with status: %d", response.StatusCode)
	}

	return nil
}
