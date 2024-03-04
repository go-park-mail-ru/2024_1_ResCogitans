package sight

import (
	"context"

	"github.com/go-park-mail-ru/2024_1_ResCogitans/internal/models/sight"
	"github.com/go-park-mail-ru/2024_1_ResCogitans/utils/response"
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

func (h *GetSights) ServeHTTP(ctx context.Context, s sight.Sight) (response.Response, error) {
	sights := sight.GetSightsList()

	return response.OK(sights), nil
}
