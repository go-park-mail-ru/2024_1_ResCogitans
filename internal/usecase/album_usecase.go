package usecase

import (
	"github.com/go-park-mail-ru/2024_1_ResCogitans/internal/entities"
	storage "github.com/go-park-mail-ru/2024_1_ResCogitans/internal/storage/storage_interfaces"
)

type AlbumUseCaseInterface interface {
	CreateAlbum(album entities.Album) (entities.Album, error)
	DeleteAlbum(album entities.Album) (entities.Album, error)
	GetAlbums(userID int) (entities.Albums, error)
}

type AlbumUseCase struct {
	AlbumStorage storage.AlbumStorageInterface
}

func NewAlbumUseCase(storage storage.AlbumStorageInterface) AlbumUseCaseInterface {
	return &AlbumUseCase{
		AlbumStorage: storage,
	}
}

func (au *AlbumUseCase) CreateAlbum(album entities.Album) (entities.Album, error) {
	row, err := au.AlbumStorage.CreateAlbum(album)
	return row, err
}

func (au *AlbumUseCase) DeleteAlbum(album entities.Album) (entities.Album, error) {
	err := au.AlbumStorage.DeleteAlbum(album.ID)
	if err != nil {
		return entities.Album{}, err
	}
	return entities.Album{ID: album.ID}, nil
}

func (au *AlbumUseCase) GetAlbums(userID int) (entities.Albums, error) {
	albums, err := au.AlbumStorage.GetAlbums(userID)
	if err != nil {
		return entities.Albums{}, err
	}
	return albums, nil
}
