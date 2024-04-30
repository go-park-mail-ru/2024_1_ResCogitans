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
// @Summary Запрос достопримечательностей
// @Description Возвращает все достопримечательности
// @Tags Достопримечательности
// @Accept json
// @Produce json
// @Success 200 {object} entities.Sights
// @Router /sights [get]
func (h *SightHandler) GetSights(_ context.Context, _ entities.SightRequest) (entities.Sights, error) {
	sights, err := h.SightUseCase.GetSightsList()
	if err != nil {
		return entities.Sights{}, err
	}
	return entities.Sights{Sight: sights}, nil
}

// GetSight godoc
// @Summary Запрос достопримечательности по id
// @Tags Достопримечательности
// @Accept json
// @Produce json
// @Success 200 {object} entities.SightComments
// @Failure 404 "Not found"
// @Router /sight/{id} [get]
func (h *SightHandler) GetSight(ctx context.Context, _ entities.SightRequest) (entities.SightComments, error) {
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

// SearchSights godoc
// @Summary Поиск достопримечательностей
// @Description Находит все достопримечательности, которые включают в себя передаваемую строку
// @Tags Достопримечательности
// @Accept json
// @Produce json
// @Param user body entities.SightRequest true "Строка для поиска"
// @Success 200 {object} entities.Sights
// @Failure 500 {object} httperrors.HttpError
// @Router /api/sight/quiz [get]
func (h *SightHandler) SearchSights(_ context.Context, requestData entities.SightRequest) (entities.Sights, error) {
	sights, err := h.SightUseCase.SearchSights(requestData.Name)
	if err != nil {
		return entities.Sights{}, err
	}
	return sights, nil
}
