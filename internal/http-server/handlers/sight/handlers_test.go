package sight_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-park-mail-ru/2024_1_ResCogitans/internal/config"
	"github.com/go-park-mail-ru/2024_1_ResCogitans/internal/models/sight"
	"github.com/go-park-mail-ru/2024_1_ResCogitans/router"
	"github.com/stretchr/testify/assert"
)

func TestGetSights(t *testing.T) {
	cfg, _ := config.LoadConfig()
	router := router.SetupRouter(cfg)
	server := httptest.NewServer(router)
	defer server.Close()

	client := server.Client()

	resp, err := client.Get(server.URL + "/sights")
	if err != nil {
		t.Fatalf("Failed to send request: %v", err)
	}
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var sights []sight.Sight
	err = json.NewDecoder(resp.Body).Decode(&sights)
	if err != nil {
		t.Fatalf("Failed to decode response body: %v", err)
	}

	assert.NotEmpty(t, sights)

	expectedFirstSight := sight.Sight{
		ID:          1,
		Rating:      4.3434,
		Name:        "Парижская башня",
		Description: "Самая высокая башня в мире.",
		City:        "Париж",
		Images:      []string{"path/to/image1.jpg", "path/to/image2.jpg"},
	}
	assert.Equal(t, expectedFirstSight, sights[0])
}
