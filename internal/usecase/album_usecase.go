package usecase

import (
	"fmt"

	"github.com/go-park-mail-ru/2024_1_ResCogitans/internal/entities"
	storage "github.com/go-park-mail-ru/2024_1_ResCogitans/internal/storage/storage_interfaces"
)

type AlbumUseCaseInterface interface {
	CreateAlbum(album entities.Album) (entities.Album, error)
	DeleteAlbum(album entities.Album) (entities.Album, error)
	GetAlbums(userID int) (entities.Albums, error)
	AddPhoto(albumID int, path string) error
	DeletePhoto(photoID int) error
	GetAlbumByID(albumID int) (entities.AlbumAndPhoto, error)
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

func (au *AlbumUseCase) AddPhoto(albumID int, path string) error {
	err := au.AlbumStorage.AddPhoto(albumID, path)
	return err
}

func (au *AlbumUseCase) DeletePhoto(photoID int) error {
	photo, err := au.AlbumStorage.DeletePhoto(photoID)
	if err != nil {
		return err
	}

	err = deleteResource(photo.Path)
	if err != nil {
		fmt.Printf("Error deleting resource: %s\n", err)
		return err
	}
	return nil
}

func (au *AlbumUseCase) GetAlbumByID(albumID int) (entities.AlbumAndPhoto, error) {
	var albumAndPhotos entities.AlbumAndPhoto

	albumInfo, err := au.AlbumStorage.GetAlbumInfo(albumID)
	if err != nil {
		return entities.AlbumAndPhoto{}, err
	}
	albumPhotos, err := au.AlbumStorage.GetAlbumPhotos(albumID)
	if err != nil {
		return entities.AlbumAndPhoto{}, err
	}

	for _, photo := range albumPhotos {
		photo.Path, err = GetDownloadLink(photo.Path)
		if err != nil {
			return entities.AlbumAndPhoto{}, err
		}
	}

	albumAndPhotos.Info = albumInfo
	albumAndPhotos.Photos = albumPhotos
	return albumAndPhotos, err
}
