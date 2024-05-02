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

func (h Album) Validate() error {
	return nil
}

func (h Albums) Validate() error {
	return nil
}
