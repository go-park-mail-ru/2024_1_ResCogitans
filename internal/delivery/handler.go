package sight

import (
	"context"

	"github.com/go-park-mail-ru/2024_1_ResCogitans/internal/entities"
	"github.com/go-park-mail-ru/2024_1_ResCogitans/internal/http-server/server/db"
	"github.com/go-park-mail-ru/2024_1_ResCogitans/utils/logger"

	sightRep "github.com/go-park-mail-ru/2024_1_ResCogitans/internal/repository/postgres"
)

// type SightUsecase struct {
// 	sightRepo sightRep.SightRepo
// }

// func (su SightUsecase) GetSights() []entities.Sight {
// 	sights := su.GetSights()
// 	return sights
// }

type SightsHandler struct{}

type Empty struct{}

type Sights struct {
	Sight []entities.Sight `json:"sights"`
}

// GetSights godoc
// @Summary Get all sights
// @Description get all sights
// @ID get-sights
// @Accept json
// @Produce json
// @Success 200 {array} sight.Sight
// @Router /sights [get]
func (h *SightsHandler) GetSights(ctx context.Context, _ entities.Sight) (Sights, error) {
	db, err := db.GetPostgres()
	if err != nil {
		logger.Logger().Error(err.Error())
	}
	sightsRepo := sightRep.NewSightRepo(db)
	sights, err := sightsRepo.GetSightsList()
	if err != nil {
		return Sights{}, err
	}

	return Sights{Sight: sights}, nil
}
