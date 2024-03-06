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
		Rating:      2.1,
		Name:        "У дяди Вани",
		Description: "Ресторан с видом на Сталинскую высотку.",
		City:        "Москва",
		Url:         "1.jpg",
	}
	assert.Equal(t, expectedFirstSight, resp.Sight[0])
}
