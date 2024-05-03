package album

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type DownloadResponse struct {
	Href string `json:"href"`
}

func getDownloadLink(filePath, token string) (DownloadResponse, error) {
	path := "jantugan/album/" + filePath
	url := "https://cloud-api.yandex.net/v1/disk/resources/download?path=" + path

	request, err := http.NewRequest("GET", url, nil)
	if err != nil {

		return DownloadResponse{}, err
	}
	request.Header.Set("Authorization", "OAuth "+token)

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return DownloadResponse{}, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return DownloadResponse{}, fmt.Errorf("download link request failed with status: %d", response.StatusCode)
	}

	var downloadResponse DownloadResponse
	if err := json.NewDecoder(response.Body).Decode(&downloadResponse); err != nil {
		return downloadResponse, err
	}

	return downloadResponse, nil
}
