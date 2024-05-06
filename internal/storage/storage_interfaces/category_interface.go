package storage

import (
	"github.com/go-park-mail-ru/2024_1_ResCogitans/internal/entities"
)

type CategoryStorageInterface interface {
	GetCategoriesList() ([]entities.Category, error)
}
