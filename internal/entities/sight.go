package entities

type Sight struct {
	ID          int     `json:"id"`
	Rating      float32 `json:"rating"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	CityID      int     `json:"cityID"`
	CountryID   int     `json:"countryID"`
	City        string  `json:"cityID"`
	Country     string  `json:"countryID"`
	Path        string  `json:"url"`
}

func (h Sight) Validate() error {
	return nil
}
