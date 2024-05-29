package sight_test

import (
	"context"
	"testing"

	"github.com/go-park-mail-ru/2024_1_ResCogitans/internal/delivery/handlers/sight"
	"github.com/go-park-mail-ru/2024_1_ResCogitans/internal/entities"
	"github.com/go-park-mail-ru/2024_1_ResCogitans/internal/mocks"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

func TestSightHandler_GetSights(t *testing.T) {
	mockSightUseCase := new(mocks.MockSightUseCase)

	handler := sight.NewSightsHandler(mockSightUseCase)

	t.Run("Successfully receiving a list of sites", func(t *testing.T) {
		// TODO подставить в sights значения
		mockSightUseCase.On("GetSightsList", context.Background()).Return([]entities.Sight{}, nil).Once()

		response, err := handler.GetSights(context.Background(), entities.Sight{})

		assert.NoError(t, err)
		assert.Equal(t, entities.Sights{Sight: []entities.Sight{}}, response)

		mockSightUseCase.AssertExpectations(t)
	})

	mockSightUseCase.Mock.ExpectedCalls = nil

	t.Run("Error getting sight", func(t *testing.T) {
		mockSightUseCase.On("GetSightsList", context.Background()).Return([]entities.Sight{}, errors.New("unexpected error")).Once()
		response, err := handler.GetSights(context.Background(), entities.Sight{})

		assert.Error(t, err)
		assert.Equal(t, entities.Sights{}, response)
		assert.Equal(t, "unexpected error", err.Error())

		mockSightUseCase.AssertExpectations(t)
	})
}

func TestSightHandler_GetSight(t *testing.T) {
	mockSightUseCase := new(mocks.MockSightUseCase)

	handler := sight.NewSightsHandler(mockSightUseCase)

	t.Run("Error getting path parameters", func(t *testing.T) {
		response, err := handler.GetSight(context.Background(), entities.Sight{})
		assert.Error(t, err)
		assert.Equal(t, "error getting id from parameters: strconv.Atoi: parsing \"\": invalid syntax", err.Error())
		assert.Equal(t, entities.SightComments{}, response)
	})

	// TODO создать запрос с параметром, где ключ - id, значение - 1
	//req, err := http.NewRequest("GET", "/api/sight/1", nil)
	//if err != nil {
	//	assert.NoError(t, err)
	//}
	//
	//params := wrapper.GetPathParams(req)
	//ctx := httputils.SetPathParamsToCtx(context.Background(), params)
	//
	//t.Run("Successfully get sight", func(t *testing.T) {
	//	mockSightUseCase""
	//})
}

func TestSightHandler_SearchSights(t *testing.T) {
	mockSightUseCase := new(mocks.MockSightUseCase)

	handler := sight.NewSightsHandler(mockSightUseCase)

	t.Run("Successfully receiving a sight", func(t *testing.T) {
		// TODO подставить в sights значения
		mockSightUseCase.On("SearchSights", context.Background(), "лес").Return(entities.Sights{}, nil).Once()
		response, err := handler.SearchSights(context.Background(), entities.Sight{Name: "лес"})

		assert.NoError(t, err)
		assert.Equal(t, entities.Sights{}, response)

		mockSightUseCase.AssertExpectations(t)
	})

	mockSightUseCase.Mock.ExpectedCalls = nil

	t.Run("Error searching sight", func(t *testing.T) {
		mockSightUseCase.On("SearchSights", context.Background(), "лес").Return(entities.Sights{}, errors.New("unexpected error")).Once()
		response, err := handler.SearchSights(context.Background(), entities.Sight{Name: "лес"})

		assert.Error(t, err)
		assert.Equal(t, entities.Sights{}, response)

		mockSightUseCase.AssertExpectations(t)
	})
}
