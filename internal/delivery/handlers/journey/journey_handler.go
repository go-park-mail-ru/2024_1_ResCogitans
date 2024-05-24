package journey

import (
	"context"
	"strconv"

	"github.com/go-park-mail-ru/2024_1_ResCogitans/internal/entities"
	"github.com/go-park-mail-ru/2024_1_ResCogitans/internal/usecase"
	"github.com/go-park-mail-ru/2024_1_ResCogitans/utils/httputils"
)

type JourneyHandler struct {
	JourneyUseCase *usecase.JourneyUseCase
}

func NewJourneyHandler(usecase *usecase.JourneyUseCase) *JourneyHandler {
	return &JourneyHandler{
		JourneyUseCase: usecase,
	}
}

func (h *JourneyHandler) CreateJourney(ctx context.Context, requestData entities.Journey) (entities.Journey, error) {
	journey, err := h.JourneyUseCase.CreateJourney(ctx, requestData)
	if err != nil {
		return entities.Journey{}, err
	}

	return journey, nil
}

func (h *JourneyHandler) EditJourney(ctx context.Context, requestData entities.Journey) (entities.Journey, error) {
	pathParams := httputils.GetPathParamsFromCtx(ctx)
	journeyID, err := strconv.Atoi(pathParams["id"])
	if err != nil {
		return entities.Journey{}, err
	}
	name := requestData.Name
	description := requestData.Description
	err = h.JourneyUseCase.EditJourney(ctx, journeyID, name, description)
	if err != nil {
		return entities.Journey{}, err
	}
	return h.JourneyUseCase.GetJourney(ctx, journeyID)
}

func (h *JourneyHandler) DeleteJourney(ctx context.Context, _ entities.Journey) (entities.Journey, error) {
	pathParams := httputils.GetPathParamsFromCtx(ctx)
	journeyID, err := strconv.Atoi(pathParams["id"])
	if err != nil {
		return entities.Journey{}, err
	}

	err = h.JourneyUseCase.DeleteJourneyByID(ctx, journeyID)
	if err != nil {
		return entities.Journey{}, err
	}

	return entities.Journey{}, nil
}

func (h *JourneyHandler) GetJourneys(ctx context.Context, _ entities.Journey) (entities.Journeys, error) {
	pathParams := httputils.GetPathParamsFromCtx(ctx)
	userID, err := strconv.Atoi(pathParams["userID"])
	if err != nil {
		return entities.Journeys{}, err
	}

	journeys, err := h.JourneyUseCase.GetJourneys(ctx, userID)
	if err != nil {
		return entities.Journeys{}, err
	}
	return entities.Journeys{Journey: journeys}, nil
}

func (h *JourneyHandler) AddJourneySight(ctx context.Context, requestData entities.JourneySightID) (entities.JourneySight, error) {
	pathParams := httputils.GetPathParamsFromCtx(ctx)
	journeyID, err := strconv.Atoi(pathParams["id"])
	if err != nil {
		return entities.JourneySight{}, err
	}

	err = h.JourneyUseCase.AddJourneySight(ctx, journeyID, requestData.ListID)
	if err != nil {
		return entities.JourneySight{}, err
	}

	return entities.JourneySight{JourneyID: journeyID}, nil
}

func (h *JourneyHandler) DeleteJourneySight(ctx context.Context, requestData entities.JourneySight) (entities.JourneySight, error) {
	pathParams := httputils.GetPathParamsFromCtx(ctx)
	journeyID, err := strconv.Atoi(pathParams["id"])
	if err != nil {
		return entities.JourneySight{}, err
	}

	err = h.JourneyUseCase.DeleteJourneySight(ctx, journeyID, requestData)
	if err != nil {
		return entities.JourneySight{}, err
	}

	return entities.JourneySight{}, nil
}

func (h *JourneyHandler) GetJourneySights(ctx context.Context, _ entities.JourneySight) (entities.JourneySights, error) {
	pathParams := httputils.GetPathParamsFromCtx(ctx)
	journeyID, err := strconv.Atoi(pathParams["id"])
	if err != nil {
		return entities.JourneySights{}, err
	}

	sights, err := h.JourneyUseCase.GetJourneySights(ctx, journeyID)
	if err != nil {
		return entities.JourneySights{}, err
	}

	journey, err := h.JourneyUseCase.GetJourney(ctx, journeyID)
	if err != nil {
		return entities.JourneySights{}, err
	}

	return entities.JourneySights{Journey: journey, Sights: sights}, nil
}
