package sight_test

import (
	"context"
	"testing"

	sight "github.com/go-park-mail-ru/2024_1_ResCogitans/internal/delivery"
	"github.com/go-park-mail-ru/2024_1_ResCogitans/internal/entities"
	"github.com/stretchr/testify/assert"
)

func TestGetSights(t *testing.T) {
	handler := &sight.SightsHandler{}

	ctx := context.Background()

	resp, err := handler.GetSights(ctx, entities.Sight{})
	if err != nil {
		t.Fatalf("Failed to get sights: %v", err)
	}

	assert.NotEmpty(t, resp.Sight)

	expectedFirstSight := entities.Sight{
		ID:          1,
		Rating:      2.1,
		Name:        "У дяди Вани",
		Description: "Ресторан с видом на Сталинскую высотку.",
		City_id:     "1",
		Country_id:  "1",
		Url:         "public/1.jpg",
	}
	assert.Equal(t, expectedFirstSight, resp.Sight[0])
}
