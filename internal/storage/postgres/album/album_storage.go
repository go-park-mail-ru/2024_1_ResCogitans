package album

import (
	"context"

	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/go-park-mail-ru/2024_1_ResCogitans/internal/entities"
	storage "github.com/go-park-mail-ru/2024_1_ResCogitans/internal/storage/storage_interfaces"
	"github.com/go-park-mail-ru/2024_1_ResCogitans/utils/logger"
	"github.com/jackc/pgx/v5/pgxpool"
)

// AlbunStorage struct
type AlbumStorage struct {
	db *pgxpool.Pool
}

// NewAlbumRepo creates sight repo
func NewAlbumStorage(db *pgxpool.Pool) storage.AlbumStorageInterface {
	return &AlbumStorage{
		db: db,
	}
}

func (as *AlbumStorage) CreateAlbum(album entities.Album) (entities.Album, error) {
	var albumID int
	ctx := context.Background()

	err := as.db.QueryRow(ctx, `INSERT INTO album(user_id, name, description) 
	VALUES($1, $2, $3) RETURNING id`, album.UserID, album.Name, album.Description).Scan(&albumID)
	if err != nil {
		logger.Logger().Error(err.Error())
		return entities.Album{}, err
	}

	return entities.Album{ID: albumID}, err
}

func (as *AlbumStorage) GetAlbums(userID int) (entities.Albums, error) {
	var albums []*entities.Album

	ctx := context.Background()

	err := pgxscan.Select(ctx, as.db, &albums, `SELECT id, name, description  
	FROM album
	WHERE user_id = $1`, userID)
	if err != nil {
		logger.Logger().Error(err.Error())
		return entities.Albums{}, err
	}

	var albumList []entities.Album
	for _, a := range albums {
		albumList = append(albumList, *a)
	}
	return entities.Albums{Albums: albumList}, nil
}

func (as *AlbumStorage) DeleteAlbum(albumID int) error {
	ctx := context.Background()

	_, err := as.db.Exec(ctx, `DELETE FROM album_photo WHERE album_id = $1`, albumID)
	if err != nil {
		logger.Logger().Error(err.Error())
		return err
	}

	_, err = as.db.Exec(ctx, `DELETE FROM album WHERE id = $1;`, albumID)
	if err != nil {
		logger.Logger().Error(err.Error())
		return err
	}

	return nil
}

func (as *AlbumStorage) AddPhoto(albumID int, path, description string) error {
	ctx := context.Background()
	_, err := as.db.Exec(ctx, `INSERT INTO album_photo(album_id, path, description) VALUES ($1, $2, $3)`, albumID, path, description)
	if err != nil {
		logger.Logger().Error(err.Error())
		return err
	}

	return nil
}

func (as *AlbumStorage) DeletePhoto(photoID int) (entities.AlbumPhoto, error) {
	var photo []*entities.AlbumPhoto
	ctx := context.Background()

	err := pgxscan.Select(ctx, as.db, &photo, `SELECT path   
	FROM album_photo
	WHERE id = $1`, photoID)
	if err != nil {
		logger.Logger().Error(err.Error())
		return entities.AlbumPhoto{}, err
	}

	_, err = as.db.Exec(ctx, `DELETE album_photo WHERE id = $1`, photoID)
	if err != nil {
		logger.Logger().Error(err.Error())
		return entities.AlbumPhoto{}, err
	}

	return *photo[0], nil
}

func (as *AlbumStorage) GetAlbumInfo(albumID int) (entities.Album, error) {
	var albums []*entities.Album

	ctx := context.Background()

	err := pgxscan.Select(ctx, as.db, &albums, `SELECT id, user_id, name, description  
	FROM album
	WHERE id = $1`, albumID)
	if err != nil {
		logger.Logger().Error(err.Error())
		return entities.Album{}, err
	}

	return *albums[0], nil
}

func (as *AlbumStorage) GetAlbumPhotos(albumID int) ([]entities.AlbumPhoto, error) {
	var albumPhotos []*entities.AlbumPhoto

	ctx := context.Background()

	err := pgxscan.Select(ctx, as.db, &albumPhotos, `SELECT id, path, description  
	FROM album_photo
	WHERE album_id = $1`, albumID)
	if err != nil {
		logger.Logger().Error(err.Error())
		return nil, err
	}

	var photoList []entities.AlbumPhoto
	for _, p := range albumPhotos {
		photoList = append(photoList, *p)
	}

	return photoList, nil
}
