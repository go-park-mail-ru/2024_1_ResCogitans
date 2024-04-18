package journey_test

import (
	"context"
	"testing"

	"github.com/go-park-mail-ru/2024_1_ResCogitans/internal/entities"
	"github.com/go-park-mail-ru/2024_1_ResCogitans/utils/wrapper"
	"github.com/stretchr/testify/assert"
)

// Get Journey
func TestGetJourney(t *testing.T) {
	handler := &JourneyHandler{}
	journeys := entities.Journey{UserID: 1}

	ctx := context.Background()

	resp, err := handler.GetJourneys(ctx, journeys)

	assert.NotEmpty(t, resp.Journey)
	assert.Nil(t, err)
}

// Get JourneySight
func TestGetJourneySights(t *testing.T) {
	handler := &JourneyHandler{}

	ctx := context.Background()
	param := make(map[string]string)
	param["id"] = "1"
	ctx = wrapper.SetPathParamsToCtx(ctx, param)

	_, err := handler.GetJourneySights(ctx, entities.JourneySight{})

	assert.Nil(t, err)
}

// Get JourneySight bad ID
func TestGetJourneySightsNotInt(t *testing.T) {
	handler := &JourneyHandler{}

	ctx := context.Background()
	param := make(map[string]string)
	param["id"] = "ok"
	ctx = wrapper.SetPathParamsToCtx(ctx, param)

	_, err := handler.GetJourneySights(ctx, entities.JourneySight{})

	assert.EqualError(t, err, "cannot parsing not integer")
}

// Create
func TestCreateJourney(t *testing.T) {
	handler := &JourneyHandler{}
	journeys := entities.Journey{
		UserID:      1,
		Name:        "siu",
		Description: "Ronaldooooo",
	}

	ctx := context.Background()

	_, err := handler.CreateJourney(ctx, journeys)

	assert.Nil(t, err)
}

// Delete
func TestDeleteJourney(t *testing.T) {
	handler := &JourneyHandler{}

	ctx := context.Background()
	param := make(map[string]string)
	param["id"] = "1"
	ctx = wrapper.SetPathParamsToCtx(ctx, param)

	_, err := handler.DeleteJourney(ctx, entities.Journey{})

	assert.Nil(t, err)
}

func TestDeleteJourneyNotInt(t *testing.T) {
	handler := &JourneyHandler{}

	ctx := context.Background()
	param := make(map[string]string)
	param["id"] = "ok"
	ctx = wrapper.SetPathParamsToCtx(ctx, param)

	_, err := handler.DeleteJourney(ctx, entities.Journey{})

	assert.EqualError(t, err, "cannot parsing not integer")
}

// Add JourneySight
func TestAddJourneySight(t *testing.T) {
	handler := &JourneyHandler{}
	journeySight := entities.JourneySight{
		SightID: 6,
	}

	ctx := context.Background()
	param := make(map[string]string)
	param["id"] = "6"
	ctx = wrapper.SetPathParamsToCtx(ctx, param)

	_, err := handler.AddJourneySight(ctx, journeySight)

	assert.Nil(t, err)
}

func TestAddJourneySightBad(t *testing.T) {
	handler := &JourneyHandler{}
	journeySight := entities.JourneySight{
		SightID: 6,
	}

	ctx := context.Background()
	param := make(map[string]string)
	param["id"] = "1"
	ctx = wrapper.SetPathParamsToCtx(ctx, param)

	_, err := handler.AddJourneySight(ctx, journeySight)

	assert.EqualError(t, err, "failed adding journey sight")
}

func TestAddJourneySightNotInt(t *testing.T) {
	handler := &JourneyHandler{}
	journeySight := entities.JourneySight{
		SightID: 6,
	}

	ctx := context.Background()
	param := make(map[string]string)
	param["id"] = "ok"
	ctx = wrapper.SetPathParamsToCtx(ctx, param)

	_, err := handler.AddJourneySight(ctx, journeySight)

	assert.EqualError(t, err, "cannot parsing not integer")
}

// Delete JourneySight
func TestDeleteJourneySight(t *testing.T) {
	handler := &JourneyHandler{}
	journeySight := entities.JourneySight{
		SightID: 6,
	}

	ctx := context.Background()
	param := make(map[string]string)
	param["id"] = "6"
	ctx = wrapper.SetPathParamsToCtx(ctx, param)

	_, err := handler.DeleteJourneySight(ctx, journeySight)

	assert.Nil(t, err)
}

func TestDeleteJourneySightNotInt(t *testing.T) {
	handler := &JourneyHandler{}
	journeySight := entities.JourneySight{
		SightID: 6,
	}

	ctx := context.Background()
	param := make(map[string]string)
	param["id"] = "ok"
	ctx = wrapper.SetPathParamsToCtx(ctx, param)

	_, err := handler.DeleteJourneySight(ctx, journeySight)

	assert.EqualError(t, err, "cannot parsing not integer")
}
