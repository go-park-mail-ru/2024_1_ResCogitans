package entities

type Sight struct {
	ID          int     `json:"id"`
	Rating      float32 `json:"rating"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	City_id     string  `json:"city"`
	Country_id  string  `json:"country"`
	Url         string  `json:"url"`
}

func (h Sight) Validate() error {
	return nil
}
