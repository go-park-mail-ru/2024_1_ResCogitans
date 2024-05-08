package usecase

import (
	"github.com/go-park-mail-ru/2024_1_ResCogitans/internal/entities"
	storage "github.com/go-park-mail-ru/2024_1_ResCogitans/internal/storage/storage_interfaces"
)

type CategoryUseCaseInterface interface {
	GetCategoriesList() ([]entities.Category, error)
}

type CategoryUseCase struct {
	CategoryStorage storage.CategoryStorageInterface
}

func NewCategoryUseCase(storage storage.CategoryStorageInterface) CategoryUseCaseInterface {
	return &CategoryUseCase{
		CategoryStorage: storage,
	}
}

func (cu *CategoryUseCase) GetCategoriesList() ([]entities.Category, error) {
	return cu.CategoryStorage.GetCategoriesList()
}
