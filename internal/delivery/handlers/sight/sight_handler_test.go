package sight_test

import (
	"context"
	"testing"

	sight "github.com/go-park-mail-ru/2024_1_ResCogitans/internal/delivery/handlers/sight"
	"github.com/go-park-mail-ru/2024_1_ResCogitans/internal/entities"
	utils "github.com/go-park-mail-ru/2024_1_ResCogitans/utils/httputils"
	"github.com/stretchr/testify/assert"
)

func TestGetSights(t *testing.T) {
	h := &sight.SightHandler{}
	h = sight.NewSightsHandler()

	ctx := context.Background()

	resp, err := h.GetSights(ctx, entities.Sight{})
	if err != nil {
		t.Fatalf("Failed to get sights: %v", err)
	}

	expectedFirstSight := entities.Sight{
		ID:          1,
		Name:        "У дяди Вани",
		Description: "Ресторан с видом на Сталинскую высотку.",
		CityID:      1,
		CountryID:   1,
		Path:        "public/1.jpg",
	}
	assert.Equal(t, expectedFirstSight, resp.Sight[0])
}

func TestGetSightsByID(t *testing.T) {
	comm := entities.Comments{}
	handler := &sight.SightHandler{}

	comm.Validate()

	ctx := context.Background()
	param := make(map[string]string)
	param["id"] = "1"
	ctx = utils.SetPathParamsToCtx(ctx, param)

	resp, err := handler.GetSight(ctx, entities.Sight{})
	if err != nil {
		t.Fatalf("Failed to get sights: %v", err)
	}

	assert.NotEmpty(t, resp.Sight)

	expectedSight := entities.Sight{
		ID:          1,
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
	handler := &sight.SightHandler{}

	ctx := context.Background()
	param := make(map[string]string)
	param["id"] = "id"
	ctx = utils.SetPathParamsToCtx(ctx, param)

	_, err := handler.GetSight(ctx, entities.Sight{})

	assert.NotNil(t, err)
}
