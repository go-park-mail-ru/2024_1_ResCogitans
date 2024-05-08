package category

import (
	"context"

	"github.com/go-park-mail-ru/2024_1_ResCogitans/internal/entities"
	storage "github.com/go-park-mail-ru/2024_1_ResCogitans/internal/storage/storage_interfaces"
	"github.com/go-park-mail-ru/2024_1_ResCogitans/utils/logger"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/georgysavva/scany/v2/pgxscan"
)

// CategoryStorage struct
type CategoryStorage struct {
	db *pgxpool.Pool
}

func NewCategoryStorage(db *pgxpool.Pool) storage.CategoryStorageInterface {
	return &CategoryStorage{
		db: db,
	}
}

func (cs *CategoryStorage) GetCategoriesList() ([]entities.Category, error) {
	var categories []*entities.Category
	ctx := context.Background()

	query := `SELECT id, name FROM category`

	err := pgxscan.Select(ctx, cs.db, &categories, query)
	if err != nil {
		logger.Logger().Error(err.Error())
		return nil, err
	}

	var categoryList []entities.Category
	for _, s := range categories {
		categoryList = append(categoryList, *s)
	}
	return categoryList, nil
}
