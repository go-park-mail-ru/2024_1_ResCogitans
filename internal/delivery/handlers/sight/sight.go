package sight

import (
	"context"
	"strconv"

	"github.com/go-park-mail-ru/2024_1_ResCogitans/internal/entities"
	"github.com/go-park-mail-ru/2024_1_ResCogitans/internal/usecase"
	"github.com/go-park-mail-ru/2024_1_ResCogitans/utils/httputils"
)

type SightHandler struct {
	SightUseCase usecase.SightUseCaseInterface
}

func NewSightsHandler(usecase usecase.SightUseCaseInterface) *SightHandler {
	return &SightHandler{
		SightUseCase: usecase,
	}
}

// GetSights godoc
// @Summary Get all sights
// @Description get all sights
// @ID get-sights
// @Accept json
// @Produce json
// @Success 200 {array} sight.Sight
// @Router /sights [get]
func (h *SightHandler) GetSights(ctx context.Context, _ entities.Sight) (entities.Sights, error) {
	params := httputils.GetQueryParamsFromCtx(ctx)
	categoryID, _ := strconv.Atoi(params["category_id"])
	sights, err := h.SightUseCase.GetSightsList(categoryID)
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
func (h *SightHandler) GetSight(ctx context.Context, _ entities.Sight) (entities.SightComments, error) {
	pathParams := httputils.GetPathParamsFromCtx(ctx)
	id, err := strconv.Atoi(pathParams["id"])
	if err != nil {
		return entities.SightComments{}, err
	}
	sight, err := h.SightUseCase.GetSightByID(id)
	if err != nil {
		return entities.SightComments{}, err
	}
	comments, err := h.SightUseCase.GetCommentsBySightID(id)
	if err != nil {
		return entities.SightComments{}, err
	}
	return entities.SightComments{Sight: sight, Comms: comments}, nil
}

func (h *SightHandler) SearchSights(_ context.Context, requestData entities.Sight) (entities.Sights, error) {
	sights, err := h.SightUseCase.SearchSights(requestData.Name)
	if err != nil {
		return entities.Sights{}, err
	}
	return sights, nil
}
