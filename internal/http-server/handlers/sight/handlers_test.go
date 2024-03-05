package sight_test

import (
	"context"
	"testing"

	"github.com/go-park-mail-ru/2024_1_ResCogitans/internal/entities"
	"github.com/go-park-mail-ru/2024_1_ResCogitans/internal/http-server/handlers/sight"
	"github.com/stretchr/testify/assert"
)

func TestGetSights(t *testing.T) {
	handler := &sight.SightsHandler{}

	ctx := context.Background()

	resp, err := handler.GetSights(ctx)
	if err != nil {
		t.Fatalf("Failed to get sights: %v", err)
	}

	assert.NotEmpty(t, resp.Sight)

	expectedFirstSight := entities.Sight{
		ID:          1,
		Rating:      4.3434,
		Name:        "Парижская башня",
		Description: "Самая высокая башня в мире.",
		City:        "Париж",
		Images:      []string{"path/to/image1.jpg", "path/to/image2.jpg"},
	}
	assert.Equal(t, expectedFirstSight, resp.Sight[0])
}
