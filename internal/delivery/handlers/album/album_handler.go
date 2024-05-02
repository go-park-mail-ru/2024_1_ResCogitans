package album

import (
	"context"

	"github.com/go-park-mail-ru/2024_1_ResCogitans/internal/entities"
	"github.com/go-park-mail-ru/2024_1_ResCogitans/internal/usecase"
	"github.com/go-park-mail-ru/2024_1_ResCogitans/utils/httputils"
)

type AlbumHandler struct {
	AlbumUseCase usecase.AlbumUseCaseInterface
}

func NewAlbumHandler(usecase usecase.AlbumUseCaseInterface) *AlbumHandler {
	return &AlbumHandler{
		AlbumUseCase: usecase,
	}
}

// CreateAlbum
// @Summary Создание нового альбома для пользователя
// @Tags Фотографии
// @Accept  json
// @Produce  json
// @Param album body entities.Album true "Данные альбома"
// @Success 200 {object} entities.Album "Успешное создание альбома"
// @Failure 500 {object} httperrors.HttpError
// @Router /api/profile/{id}/album/create [post]
func (h *AlbumHandler) CreateAlbum(ctx context.Context, album entities.Album) (entities.Album, error) {
	userID, err := httputils.GetUserFromCtx(ctx)
	if err != nil {
		return entities.Album{}, err
	}

	album.UserID = userID

	album, err = h.AlbumUseCase.CreateAlbum(album)
	if err != nil {
		return entities.Album{}, err
	}
	return album, nil
}

// DeleteAlbum
// @Summary Удаление альбома пользователя
// @Tags Фотографии
// @Accept  json
// @Produce  json
// @Param album body entities.Album true "Данные альбома"
// @Success 200 {object} entities.Album "Успешное удаление альбома"
// @Failure 500 {object} httperrors.HttpError
// @Router /api/profile/{id}/album/delete [post]
func (h *AlbumHandler) DeleteAlbum(_ context.Context, album entities.Album) (entities.Album, error) {
	_, err := h.AlbumUseCase.DeleteAlbum(album)
	if err != nil {
		return entities.Album{}, err
	}
	return entities.Album{}, nil
}

// GetAlbums
// @Summary Получение альбома пользователя
// @Tags Фотографии
// @Accept  json
// @Produce  json
// @Success 200 {object} entities.Album "Успешное удаление альбома"
// @Failure 500 {object} httperrors.HttpError
// @Router /api/profile/{id}/albums [get]
func (h *AlbumHandler) GetAlbums(ctx context.Context, _ entities.Album) (entities.Albums, error) {
	userID, err := httputils.GetUserFromCtx(ctx)
	if err != nil {
		return entities.Albums{}, err
	}

	album, err := h.AlbumUseCase.GetAlbums(userID)
	if err != nil {
		return entities.Albums{}, err
	}
	return album, nil
}
