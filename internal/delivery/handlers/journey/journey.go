package journey

import (
	"context"
	"strconv"

	"github.com/go-park-mail-ru/2024_1_ResCogitans/internal/entities"
	"github.com/go-park-mail-ru/2024_1_ResCogitans/internal/usecase"
	"github.com/go-park-mail-ru/2024_1_ResCogitans/utils/httputils"
	"github.com/go-park-mail-ru/2024_1_ResCogitans/utils/logger"
)

type JourneyHandler struct {
	JourneyUseCase usecase.JourneyUseCaseInterface
}

func NewJourneyHandler(usecase usecase.JourneyUseCaseInterface) *JourneyHandler {
	return &JourneyHandler{
		JourneyUseCase: usecase,
	}
}

func (h *JourneyHandler) CreateJourney(_ context.Context, requestData entities.Journey) (entities.Journey, error) {
	dataStr := make(map[string]string)
	dataInt := make(map[string]int)

	dataInt["userID"] = requestData.UserID
	dataStr["name"] = requestData.Name
	dataStr["description"] = requestData.Description

	journey, err := h.JourneyUseCase.CreateJourney(dataInt, dataStr)
	if err != nil {
		return entities.Journey{}, err
	}

	return journey, nil
}

func (h *JourneyHandler) DeleteJourney(ctx context.Context, _ entities.Journey) (entities.Journey, error) {
	pathParams := httputils.GetPathParamsFromCtx(ctx)
	journeyID, err := strconv.Atoi(pathParams["id"])
	if err != nil {
		logger.Logger().Error("Cannot convert string to integer to get sight")
		return entities.Journey{}, err
	}

	dataInt := make(map[string]int)
	dataInt["journeyID"] = journeyID

	err = h.JourneyUseCase.DeleteJourneyByID(dataInt)
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

	journeys, err := h.JourneyUseCase.GetJourneys(userID)
	if err != nil {
		return entities.Journeys{}, err
	}
	return entities.Journeys{Journey: journeys}, nil
}

func (h *JourneyHandler) AddJourneySight(ctx context.Context, requestData entities.JourneySight) (entities.JourneySight, error) {
	pathParams := httputils.GetPathParamsFromCtx(ctx)
	journeyID, err := strconv.Atoi(pathParams["id"])
	if err != nil {
		return entities.JourneySight{}, err
	}

	dataInt := make(map[string]int)
	dataInt["journeyID"] = journeyID
	dataInt["sightID"] = requestData.SightID

	err = h.JourneyUseCase.AddJourneySight(dataInt)
	if err != nil {
		return entities.JourneySight{}, err
	}

	return entities.JourneySight{}, nil
}

func (h *JourneyHandler) DeleteJourneySight(ctx context.Context, requestData entities.JourneySight) (entities.JourneySight, error) {
	pathParams := httputils.GetPathParamsFromCtx(ctx)
	journeyID, err := strconv.Atoi(pathParams["id"])
	if err != nil {
		logger.Logger().Error("Cannot convert string to integer to get sight")
		return entities.JourneySight{}, err
	}

	dataInt := make(map[string]int)
	dataInt["journeyID"] = journeyID
	dataInt["sightID"] = requestData.SightID

	err = h.JourneyUseCase.DeleteJourneySight(dataInt)
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

	sights, err := h.JourneyUseCase.GetJourneySights(journeyID)

	if err != nil {
		return entities.JourneySights{}, err
	}

	journey, err := h.JourneyUseCase.GetJourney(journeyID)
	if err != nil {
		return entities.JourneySights{}, err
	}

	return entities.JourneySights{Journey: journey, Sight: sights}, nil
}
