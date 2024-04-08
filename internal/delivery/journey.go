package delivery

import (
	"context"
	"errors"
	"strconv"

	"github.com/go-park-mail-ru/2024_1_ResCogitans/internal/entities"
	"github.com/go-park-mail-ru/2024_1_ResCogitans/internal/http-server/server/db"
	sightRep "github.com/go-park-mail-ru/2024_1_ResCogitans/internal/repository/postgres"
	"github.com/go-park-mail-ru/2024_1_ResCogitans/utils/logger"
	"github.com/go-park-mail-ru/2024_1_ResCogitans/utils/wrapper"
)

type JourneyHandler struct{}

func (h JourneyHandler) CreateJourney(ctx context.Context, requestData entities.Journey) (entities.Journey, error) {
	db, err := db.GetPostgres()

	if err != nil {
		logger.Logger().Error(err.Error())
	}

	dataStr := make(map[string]string)
	dataInt := make(map[string]int)

	dataInt["userID"] = requestData.UserID
	dataStr["name"] = requestData.Name
	dataStr["description"] = requestData.Description

	sightsRepo := sightRep.NewSightRepo(db)
	err = sightsRepo.CreateJourney(dataInt, dataStr)

	if err != nil {
		return entities.Journey{}, errors.New("cannot create journey")
	}

	return entities.Journey{}, nil
}

func (h JourneyHandler) DeleteJourney(ctx context.Context, requestData entities.Journey) (entities.Journey, error) {
	db, err := db.GetPostgres()

	if err != nil {
		logger.Logger().Error(err.Error())
	}

	pathParams := wrapper.GetPathParamsFromCtx(ctx)
	journeyID, err := strconv.Atoi(pathParams["id"])
	if err != nil {
		logger.Logger().Error("Cannot convert string to integer to get sight")
		return entities.Journey{}, err
	}

	dataInt := make(map[string]int)

	dataInt["journeyID"] = journeyID

	sightsRepo := sightRep.NewSightRepo(db)
	err = sightsRepo.DeleteJourneyByID(dataInt)

	if err != nil {
		return entities.Journey{}, errors.New("cannot delete journey")
	}

	return entities.Journey{}, nil
}

func (h *JourneyHandler) GetJourneys(ctx context.Context, requestData entities.Journey) (entities.Journeys, error) {
	db, err := db.GetPostgres()
	if err != nil {
		logger.Logger().Error(err.Error())
	}

	sightsRepo := sightRep.NewSightRepo(db)
	journeys, _ := sightsRepo.GetJourneys(requestData.UserID)

	return entities.Journeys{Journey: journeys}, err
}
