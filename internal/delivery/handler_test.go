package sight_test

import (
	"context"
	"testing"

	sight "github.com/go-park-mail-ru/2024_1_ResCogitans/internal/delivery"
	"github.com/go-park-mail-ru/2024_1_ResCogitans/internal/entities"
	"github.com/go-park-mail-ru/2024_1_ResCogitans/utils/wrapper"
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
		CityID:      1,
		CountryID:   1,
		Path:        "public/1.jpg",
	}
	assert.Equal(t, expectedFirstSight, resp.Sight[0])
}

func TestGetSightsByID(t *testing.T) {
	comm := sight.Comments{}
	handler := &sight.SightsHandler{}

	comm.Validate()

	ctx := context.Background()
	param := make(map[string]string)
	param["id"] = "1"
	ctx = wrapper.SetPathParamsToCtx(ctx, param)

	resp, err := handler.GetSightByID(ctx, entities.Sight{})
	if err != nil {
		t.Fatalf("Failed to get sights: %v", err)
	}

	assert.NotEmpty(t, resp.Sight)

	expectedSight := entities.Sight{
		ID:          1,
		Rating:      2.1,
		Name:        "У дяди Вани",
		Description: "Ресторан с видом на Сталинскую высотку.",
		CityID:      1,
		CountryID:   1,
		City:        "Москва",
		Country:     "Россия",
		Path:        "public/1.jpg",
	}
	assert.Equal(t, expectedSight, resp.Sight)
}

func TestGetSightsByIDNotInt(t *testing.T) {
	handler := &sight.SightsHandler{}

	ctx := context.Background()
	param := make(map[string]string)
	param["id"] = "id"
	ctx = wrapper.SetPathParamsToCtx(ctx, param)

	_, err := handler.GetSightByID(ctx, entities.Sight{})

	assert.NotNil(t, err)
}
