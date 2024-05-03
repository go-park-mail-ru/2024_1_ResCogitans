package entities

type Album struct {
	ID          int    `json:"albumID"`
	UserID      int    `json:"userID"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type Albums struct {
	Albums []Album `json:"albums"`
}

type AlbumPhoto struct {
	ID          int    `json:"photoID`
	AlbumID     int    `json:"albumID"`
	Path        string `json:"path"`
	Description string `json:"description"`
}

func (h Album) Validate() error {
	return nil
}

func (h Albums) Validate() error {
	return nil
}

func (h AlbumPhoto) Validate() error {
	return nil
}
