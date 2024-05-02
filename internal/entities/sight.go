package entities

type SightIDRequest struct {
	ID string `json:"sightID"`
}

func (h SightIDRequest) Validate() error {
	return nil
}

type SightRequest struct {
	Name string `json:"name"`
}

func (h SightRequest) Validate() error {
	return nil
}

type SightsList struct {
	ListID []int `json:"sightIDs"`
}

func (h SightsList) Validate() error {
	return nil
}

type Sight struct {
	ID          int      `json:"id"`
	Rating      *float64 `json:"rating"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	CityID      int      `json:"cityID"`
	CountryID   int      `json:"countryID"`
	City        string   `json:"city"`
	Country     string   `json:"country"`
	Longitude   *float64 `json:"longitude"`
	Latitude    *float64 `json:"latitude"`
	Path        string   `json:"url"`
}

type Sights struct {
	Sight []Sight `json:"sights"`
}

type SightComments struct {
	Sight Sight     `json:"sight"`
	Comms []Comment `json:"comments"`
}
