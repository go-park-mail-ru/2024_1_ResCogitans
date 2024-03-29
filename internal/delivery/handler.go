package sight

import (
	"context"

	"github.com/go-park-mail-ru/2024_1_ResCogitans/internal/entities"
	"github.com/go-park-mail-ru/2024_1_ResCogitans/internal/http-server/server/db"

	sightRep "github.com/go-park-mail-ru/2024_1_ResCogitans/internal/repository/postgres"
)

// GetSights godoc
// @Summary Get all sights
// @Description get all sights
// @ID get-sights
// @Accept json
// @Produce json
// @Success 200 {array} sight.Sight
// @Router /sights [get]
type SightUsecase struct {
	sightRepo sightRep.SightRepo
}

func (su SightUsecase) GetSights() []entities.Sight {
	sights := su.GetSights()
	return sights
}

type SightsHandler struct{}

type Sights struct {
	Sight []entities.Sight `json:"sights"`
}

func (h *SightsHandler) GetSights(ctx context.Context, _ entities.Sight) (Sights, error) {
	sightsRepo := sightRep.NewSightRepo(db.GetPostgres())
	sights, err := sightsRepo.GetSightsList()
	if err != nil {
		return Sights{}, err
	}

	var sightList []entities.Sight
	for _, s := range sights {
		sightList = append(sightList, *s)
	}

	return Sights{Sight: sightList}, nil
}
