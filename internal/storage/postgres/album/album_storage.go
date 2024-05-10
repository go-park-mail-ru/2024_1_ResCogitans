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

	query := `INSERT INTO album(user_id, name, description) 
	VALUES($1, $2, $3) RETURNING id`

	err := as.db.QueryRow(ctx, query, album.UserID, album.Name, album.Description).Scan(&albumID)
	if err != nil {
		logger.Logger().Error(err.Error())
		return entities.Album{}, err
	}

	return entities.Album{ID: albumID}, err
}

func (as *AlbumStorage) GetAlbums(userID int) (entities.Albums, error) {
	var albums []entities.Album
	ctx := context.Background()

	query := `SELECT id, name, description  
	FROM album
	WHERE user_id = $1`

	err := pgxscan.Select(ctx, as.db, &albums, query, userID)
	if err != nil {
		logger.Logger().Error(err.Error())
		return entities.Albums{}, err
	}

	return entities.Albums{Albums: albums}, nil
}

func (as *AlbumStorage) DeleteAlbum(albumID int) error {
	ctx := context.Background()

	query := `DELETE FROM album_photo WHERE album_id = $1`

	_, err := as.db.Exec(ctx, query, albumID)
	if err != nil {
		logger.Logger().Error(err.Error())
		return err
	}

	query = `DELETE FROM album WHERE id = $1;`

	_, err = as.db.Exec(ctx, query, albumID)
	if err != nil {
		logger.Logger().Error(err.Error())
		return err
	}

	return nil
}

func (as *AlbumStorage) AddPhoto(albumID int, path, description string) error {
	ctx := context.Background()

	query := `INSERT INTO album_photo(album_id, path, description) VALUES ($1, $2, $3)`

	_, err := as.db.Exec(ctx, query, albumID, path, description)
	if err != nil {
		logger.Logger().Error(err.Error())
		return err
	}

	return nil
}

func (as *AlbumStorage) DeletePhoto(photoID int) (entities.AlbumPhoto, error) {
	var photo entities.AlbumPhoto
	ctx := context.Background()

	query := `SELECT path   
	FROM album_photo
	WHERE id = $1`

	err := pgxscan.Get(ctx, as.db, &photo, query, photoID)
	if err != nil {
		logger.Logger().Error(err.Error())
		return entities.AlbumPhoto{}, err
	}

	query = `DELETE FROM album_photo WHERE id = $1`

	_, err = as.db.Exec(ctx, query, photoID)
	if err != nil {
		logger.Logger().Error(err.Error())
		return entities.AlbumPhoto{}, err
	}

	return photo, nil
}

func (as *AlbumStorage) GetAlbumInfo(albumID int) (entities.Album, error) {
	var album entities.Album

	ctx := context.Background()

	query := `SELECT id, user_id, name, description  
	FROM album
	WHERE id = $1`

	err := pgxscan.Get(ctx, as.db, &album, query, albumID)
	if err != nil {
		logger.Logger().Error(err.Error())
		return entities.Album{}, err
	}

	return album, nil
}

func (as *AlbumStorage) GetAlbumPhotos(albumID int) ([]entities.AlbumPhoto, error) {
	var albumPhotos []entities.AlbumPhoto

	ctx := context.Background()

	query := `SELECT 
		id, 
		path, 
		description  
	FROM album_photo
	WHERE album_id = $1`

	err := pgxscan.Select(ctx, as.db, &albumPhotos, query, albumID)
	if err != nil {
		logger.Logger().Error(err.Error())
		return nil, err
	}

	return albumPhotos, nil
}
