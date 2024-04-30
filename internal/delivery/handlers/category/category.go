package category

import (
	"context"

	"github.com/go-park-mail-ru/2024_1_ResCogitans/internal/entities"
	"github.com/go-park-mail-ru/2024_1_ResCogitans/internal/usecase"
)

type CategoryHandler struct {
	CategoryUseCase usecase.CategoryUseCaseInterface
}

func NewCategoryHandler(usecase usecase.CategoryUseCaseInterface) *CategoryHandler {
	return &CategoryHandler{
		CategoryUseCase: usecase,
	}
}

func (h *CategoryHandler) GetCategories(ctx context.Context, _ entities.Category) (entities.Categories, error) {
	categories, err := h.CategoryUseCase.GetCategoriesList()
	if err != nil {
		return entities.Categories{}, err
	}
	return entities.Categories{Category: categories}, nil
}
