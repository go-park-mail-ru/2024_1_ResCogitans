package entities

type Category struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func (h Category) Validate() error {
	return nil
}

type Categories struct {
	Category []Category `json:"categories"`
}
