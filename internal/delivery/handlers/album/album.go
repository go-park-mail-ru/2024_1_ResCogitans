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

// CreateAlbum godoc
// @Summary Create new album
// @Description create new album
// @ID CreateAlbum
// @Accept json
// @Produce json
// @Success 200 Album
// @Router /sight [get]
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

// GetSights godoc
// @Summary Get all sights
// @Description get all sights
// @ID get-sights
// @Accept json
// @Produce json
// @Success 200 {array} sight.Sight
// @Router /sights [get]
func (h *AlbumHandler) DeleteAlbum(_ context.Context, album entities.Album) (entities.Album, error) {
	_, err := h.AlbumUseCase.DeleteAlbum(album)
	if err != nil {
		return entities.Album{}, err
	}
	return entities.Album{}, nil
}

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
