package sight

import (
	"context"
	"net/http"

	"github.com/go-park-mail-ru/2024_1_ResCogitans/internal/models/sight"
)

// GetSights godoc
// @Summary Get all sights
// @Description get all sights
// @ID get-sights
// @Accept json
// @Produce json
// @Success 200 {array} sight.Sight
// @Router /sights [get]
type GetSights struct{}

func (h *GetSights) ServeHTTP(ctx context.Context, w http.ResponseWriter, r *http.Request) (interface{}, error) {
	sights := sight.GetSightsList()

	return sights, nil
}
