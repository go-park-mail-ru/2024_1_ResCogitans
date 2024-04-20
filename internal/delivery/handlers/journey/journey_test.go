package journey

import (
	"context"
	"testing"

	journey "github.com/go-park-mail-ru/2024_1_ResCogitans/internal/delivery"
	"github.com/go-park-mail-ru/2024_1_ResCogitans/internal/entities"
	"github.com/go-park-mail-ru/2024_1_ResCogitans/utils/wrapper"
	"github.com/stretchr/testify/assert"
)

// Get Journey
func TestGetJourney(t *testing.T) {
	handler := &journey.JourneyHandler{}
	journeys := entities.Journey{UserID: 1}

	ctx := context.Background()

	resp, err := handler.GetJourneys(ctx, journeys)

	assert.NotEmpty(t, resp.Journey)
	assert.Nil(t, err)
}

// Get JourneySight
func TestGetJourneySights(t *testing.T) {
	handler := &journey.JourneyHandler{}

	ctx := context.Background()
	param := make(map[string]string)
	param["id"] = "1"
	ctx = wrapper.SetPathParamsToCtx(ctx, param)

	_, err := handler.GetJourneySights(ctx, entities.JourneySight{})

	assert.Nil(t, err)
}

// Get JourneySight bad ID
func TestGetJourneySightsNotInt(t *testing.T) {
	handler := &journey.JourneyHandler{}

	ctx := context.Background()
	param := make(map[string]string)
	param["id"] = "ok"
	ctx = wrapper.SetPathParamsToCtx(ctx, param)

	_, err := handler.GetJourneySights(ctx, entities.JourneySight{})

	assert.EqualError(t, err, "cannot parsing not integer")
}

// Create
func TestCreateJourney(t *testing.T) {
	handler := &journey.JourneyHandler{}
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
	handler := &journey.JourneyHandler{}

	ctx := context.Background()
	param := make(map[string]string)
	param["id"] = "1"
	ctx = wrapper.SetPathParamsToCtx(ctx, param)

	_, err := handler.DeleteJourney(ctx, entities.Journey{})

	assert.Nil(t, err)
}

func TestDeleteJourneyNotInt(t *testing.T) {
	handler := &journey.JourneyHandler{}

	ctx := context.Background()
	param := make(map[string]string)
	param["id"] = "ok"
	ctx = wrapper.SetPathParamsToCtx(ctx, param)

	_, err := handler.DeleteJourney(ctx, entities.Journey{})

	assert.EqualError(t, err, "cannot parsing not integer")
}

// Add JourneySight
func TestAddJourneySight(t *testing.T) {
	ids := []int{1, 2, 3, 4}
	handler := &journey.JourneyHandler{}
	JourneySightID := entities.JourneySightID{
		ListID: ids,
	}

	ctx := context.Background()
	param := make(map[string]string)
	param["id"] = "6"
	ctx = wrapper.SetPathParamsToCtx(ctx, param)

	_, err := handler.AddJourneySight(ctx, JourneySightID)

	assert.Nil(t, err)
}

func TestAddJourneySightBad(t *testing.T) {
	handler := &journey.JourneyHandler{}
	ids := []int{1, 2, 3, 13, 18}
	JourneySightID := entities.JourneySightID{
		ListID: ids,
	}

	ctx := context.Background()
	param := make(map[string]string)
	param["id"] = "1"
	ctx = wrapper.SetPathParamsToCtx(ctx, param)

	_, err := handler.AddJourneySight(ctx, JourneySightID)

	assert.EqualError(t, err, "failed adding journey sight")
}

func TestAddJourneySightNotInt(t *testing.T) {
	handler := &journey.JourneyHandler{}

	ctx := context.Background()
	param := make(map[string]string)
	param["id"] = "ok"
	ctx = wrapper.SetPathParamsToCtx(ctx, param)

	_, err := handler.AddJourneySight(ctx, entities.JourneySightID{})

	assert.EqualError(t, err, "cannot parsing not integer")
}

// Delete JourneySight
func TestDeleteJourneySight(t *testing.T) {
	handler := &journey.JourneyHandler{}
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
	handler := &journey.JourneyHandler{}
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
