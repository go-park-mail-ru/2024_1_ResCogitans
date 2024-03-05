package sight

import (
	"context"

	"github.com/go-park-mail-ru/2024_1_ResCogitans/internal/entities"
)

// GetSights godoc
// @Summary Get all sights
// @Description get all sights
// @ID get-sights
// @Accept json
// @Produce json
// @Success 200 {array} sight.Sight
// @Router /sights [get]
type SightsHandler struct{}

type Sights struct {
	Sight []entities.Sight `json:"sights"`
}

func (h *SightsHandler) GetSights(ctx context.Context) (Sights, error) {
	sights := entities.GetSightsList()

	return Sights{Sight: sights}, nil
}
