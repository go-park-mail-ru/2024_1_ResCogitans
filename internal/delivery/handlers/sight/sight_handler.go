package sight

import (
	"context"
	"strconv"

	"github.com/go-park-mail-ru/2024_1_ResCogitans/internal/entities"
	"github.com/go-park-mail-ru/2024_1_ResCogitans/internal/usecase"
	"github.com/go-park-mail-ru/2024_1_ResCogitans/utils/httputils"
)

type SightHandler struct {
	SightUseCase *usecase.SightUseCase
}

func NewSightsHandler(usecase *usecase.SightUseCase) *SightHandler {
	return &SightHandler{
		SightUseCase: usecase,
	}
}

func (h *SightHandler) GetSights(ctx context.Context, _ entities.Sight) (entities.Sights, error) {
	sights, err := h.SightUseCase.GetSightsList(ctx)
	if err != nil {
		return entities.Sights{}, err
	}
	return entities.Sights{Sight: sights}, nil
}

func (h *SightHandler) GetSight(ctx context.Context, _ entities.Sight) (entities.SightComments, error) {
	pathParams := httputils.GetPathParamsFromCtx(ctx)
	id, err := strconv.Atoi(pathParams["id"])
	if err != nil {
		return entities.SightComments{}, err
	}
	sight, err := h.SightUseCase.GetSightByID(ctx, id)
	if err != nil {
		return entities.SightComments{}, err
	}
	comments, err := h.SightUseCase.GetCommentsBySightID(ctx, id)
	if err != nil {
		return entities.SightComments{}, err
	}
	return entities.SightComments{Sight: sight, Comms: comments}, nil
}

func (h *SightHandler) SearchSights(ctx context.Context, requestData entities.Sight) (entities.Sights, error) {
	sights, err := h.SightUseCase.SearchSights(ctx, requestData.Name)
	if err != nil {
		return entities.Sights{}, err
	}
	return sights, nil
}
