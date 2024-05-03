package journey

import (
	"context"
	"strconv"

	"github.com/go-park-mail-ru/2024_1_ResCogitans/internal/entities"
	"github.com/go-park-mail-ru/2024_1_ResCogitans/internal/usecase"
	"github.com/go-park-mail-ru/2024_1_ResCogitans/utils/httputils"
)

type JourneyHandler struct {
	JourneyUseCase usecase.JourneyUseCaseInterface
}

func NewJourneyHandler(usecase usecase.JourneyUseCaseInterface) *JourneyHandler {
	return &JourneyHandler{
		JourneyUseCase: usecase,
	}
}

// CreateJourney godoc
// @Summary Создание поездки
// @Description Создает поездку для пользователя, названия не могут повторяться (пока что вообще у всех)
// @Tags Поездки
// @Accept json
// @Produce json
// @Param user body entities.JourneyRequest true "Данные для создания поездки"
// @Success 200 {object} entities.Journey
// @Failure 500 {object} httperrors.HttpError
// @Router /api/trip/create [post]
func (h *JourneyHandler) CreateJourney(_ context.Context, requestData entities.JourneyRequest) (entities.Journey, error) {
	journey, err := h.JourneyUseCase.CreateJourney(requestData)
	if err != nil {
		return entities.Journey{}, err
	}
	return journey, nil
}

// EditJourney godoc
// @Summary Редактирование поездки
// @Description Можно изменить название и описание поездки
// @Tags Поездки
// @Accept json
// @Produce json
// @Param user body entities.JourneyRequest true "Данные для создания поездки"
// @Success 200 {object} entities.Journey
// @Failure 500 {object} httperrors.HttpError
// @Router /api/trip/{id}/edit [post]
func (h *JourneyHandler) EditJourney(ctx context.Context, requestData entities.JourneyRequest) (entities.Journey, error) {
	pathParams := httputils.GetPathParamsFromCtx(ctx)
	journeyID, err := strconv.Atoi(pathParams["id"])
	if err != nil {
		return entities.Journey{}, err
	}

	name := requestData.Name
	description := requestData.Description
	err = h.JourneyUseCase.EditJourney(journeyID, name, description)
	if err != nil {
		return entities.Journey{}, err
	}
	return h.JourneyUseCase.GetJourney(journeyID)
}

// DeleteJourney godoc
// @Summary Удаление поездки
// @Tags Поездки
// @Accept json
// @Produce json
// @Success 200 {object} entities.Journey
// @Failure 500 {object} httperrors.HttpError
// @Router /api/trip/{id}/delete [post]
func (h *JourneyHandler) DeleteJourney(ctx context.Context, _ entities.JourneyRequest) (entities.Journey, error) {
	pathParams := httputils.GetPathParamsFromCtx(ctx)
	journeyID, err := strconv.Atoi(pathParams["id"])
	if err != nil {
		return entities.Journey{}, err
	}

	err = h.JourneyUseCase.DeleteJourneyByID(journeyID)
	if err != nil {
		return entities.Journey{}, err
	}

	return entities.Journey{}, nil
}

// GetJourneys godoc
// @Summary Получение поездок по id пользователя
// @Tags Поездки
// @Accept json
// @Produce json
// @Success 200 {object} entities.Journeys
// @Failure 500 {object} httperrors.HttpError
// @Router /api/{userID}/trips [get]
func (h *JourneyHandler) GetJourneys(ctx context.Context, _ entities.JourneyRequest) (entities.Journeys, error) {
	pathParams := httputils.GetPathParamsFromCtx(ctx)
	userID, err := strconv.Atoi(pathParams["userID"])
	if err != nil {
		return entities.Journeys{}, err
	}

	journeys, err := h.JourneyUseCase.GetJourneys(userID)
	if err != nil {
		return entities.Journeys{}, err
	}
	return entities.Journeys{Journey: journeys}, nil
}

// AddJourneySight godoc
// @Summary Добавление достопримечательностей в поездку
// @Description Принимает список айди достопримечательностей, которые необходимо добавить в поездку, берет айди поездки из url
// @Tags Поездки
// @Accept json
// @Produce json
// @Param user body entities.SightsList true "Айди поездок"
// @Success 200 {object} entities.JourneySight
// @Failure 500 {object} httperrors.HttpError
// @Router /api/trip/{id}/sight/add [post]
func (h *JourneyHandler) AddJourneySight(ctx context.Context, requestData entities.SightsList) (entities.JourneySight, error) {
	pathParams := httputils.GetPathParamsFromCtx(ctx)
	journeyID, err := strconv.Atoi(pathParams["id"])
	if err != nil {
		return entities.JourneySight{}, err
	}

	err = h.JourneyUseCase.AddJourneySight(journeyID, requestData.ListID)
	if err != nil {
		return entities.JourneySight{}, err
	}

	return entities.JourneySight{JourneyID: journeyID}, nil
}

// DeleteJourneySight godoc
// @Summary Удаление достопримечательности из поездки
// @Tags Поездки
// @Accept json
// @Produce json
// @Param user body entities.SightIDRequest true "Айди достопримечательности"
// @Success 200 {object} entities.JourneySight
// @Failure 500 {object} httperrors.HttpError
// @Router /api/trip/{id}/sight/delete [post]
func (h *JourneyHandler) DeleteJourneySight(ctx context.Context, requestData entities.SightIDRequest) (entities.JourneySight, error) {
	pathParams := httputils.GetPathParamsFromCtx(ctx)
	journeyID, err := strconv.Atoi(pathParams["id"])
	if err != nil {
		return entities.JourneySight{}, err
	}

	err = h.JourneyUseCase.DeleteJourneySight(journeyID, requestData)
	if err != nil {
		return entities.JourneySight{}, err
	}

	return entities.JourneySight{}, nil
}

// GetJourneySights godoc
// @Summary Получение достопримечательностей из поездки
// @Tags Поездки
// @Accept json
// @Produce json
// @Success 200 {object} entities.JourneySights
// @Failure 500 {object} httperrors.HttpError
// @Router /api/trip/{id} [get]
func (h *JourneyHandler) GetJourneySights(ctx context.Context, _ entities.JourneyRequest) (entities.JourneySights, error) {
	pathParams := httputils.GetPathParamsFromCtx(ctx)
	journeyID, err := strconv.Atoi(pathParams["id"])
	if err != nil {
		return entities.JourneySights{}, err
	}

	sights, err := h.JourneyUseCase.GetJourneySights(journeyID)
	if err != nil {
		return entities.JourneySights{}, err
	}

	journey, err := h.JourneyUseCase.GetJourney(journeyID)
	if err != nil {
		return entities.JourneySights{}, err
	}

	return entities.JourneySights{Journey: journey, Sights: sights}, nil
}
