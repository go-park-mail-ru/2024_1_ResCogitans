package storage

import (
	"github.com/go-park-mail-ru/2024_1_ResCogitans/internal/entities"
)

type AlbumStorageInterface interface {
	CreateAlbum(album entities.Album) (entities.Album, error)
	DeleteAlbum(albumID int) error
	GetAlbums(userID int) (entities.Albums, error)
}
