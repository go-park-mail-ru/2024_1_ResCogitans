package delivery

import (
	"context"
	"strconv"

	"github.com/go-park-mail-ru/2024_1_ResCogitans/internal/entities"
	"github.com/go-park-mail-ru/2024_1_ResCogitans/internal/http-server/server/db"
	"github.com/go-park-mail-ru/2024_1_ResCogitans/utils/logger"
	"github.com/go-park-mail-ru/2024_1_ResCogitans/utils/wrapper"

	sightRep "github.com/go-park-mail-ru/2024_1_ResCogitans/internal/repository/postgres"
)

type SightsHandler struct{}

type SightComments struct {
	Sight entities.Sight     `json:"sight"`
	Comms []entities.Comment `json:"comments"`
}

// GetSights godoc
// @Summary Get all sights
// @Description get all sights
// @ID get-sights
// @Accept json
// @Produce json
// @Success 200 {array} sight.Sight
// @Router /sights [get]
func (h *SightsHandler) GetSights(ctx context.Context, _ entities.Sight) (entities.Sights, error) {
	db, err := db.GetPostgres()

	if err != nil {
		logger.Logger().Error(err.Error())
	}
	sightsRepo := sightRep.NewSightRepo(db)
	sights, err := sightsRepo.GetSightsList()
	if err != nil {
		return entities.Sights{}, err
	}

	return entities.Sights{Sight: sights}, nil
}

// GetSights godoc
// @Summary Get sight by id
// @Description get sight by id
// @Accept json
// @Produce json
// @Success 200 SightComments
// @Failure 404 "Not found"
// @Router /sight/{id} [get]
func (h *SightsHandler) GetSightByID(ctx context.Context, requestData entities.Sight) (SightComments, error) {
	db, err := db.GetPostgres()
	if err != nil {
		logger.Logger().Error(err.Error())
	}

	pathParams := wrapper.GetPathParamsFromCtx(ctx)
	id, err := strconv.Atoi(pathParams["id"])
	if err != nil {
		logger.Logger().Error("Cannot convert string to integer to get sight")
		return SightComments{}, err
	}

	sightsRepo := sightRep.NewSightRepo(db)
	sight, _ := sightsRepo.GetSightByID(id)

	comments, err := sightsRepo.GetCommentsBySightID(id)

	return SightComments{Sight: sight, Comms: comments}, err
}
