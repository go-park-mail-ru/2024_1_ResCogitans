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

	resp, err := handler.GetSights(ctx, entities.Sight{})
	if err != nil {
		t.Fatalf("Failed to get sights: %v", err)
	}

	assert.NotEmpty(t, resp.Sight)

	expectedFirstSight := entities.Sight{
		ID:          1,
<<<<<<< HEAD
		Rating:      4.3434,
		Name:        "Парижская башня",
		Description: "Самая высокая башня в мире.",
		City:        "Париж",
		Url:         "public/1.jpg",
=======
		Rating:      2.1,
		Name:        "У дяди Вани",
		Description: "Ресторан с видом на Сталинскую высотку.",
		City:        "Москва",
		Url:         "1.jpg",
>>>>>>> auth
	}
	assert.Equal(t, expectedFirstSight, resp.Sight[0])
}
