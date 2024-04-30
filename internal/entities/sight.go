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

type Sight struct {
	ID          int
	Rating      *float32
	Name        string
	Description string
	CityID      int
	CountryID   int
	City        string
	Country     string
	Path        string
}

type SightsList struct {
	ListID []int `json:"sightIDs"`
}

func (h SightsList) Validate() error {
	return nil
}

type Sights struct {
	Sight []Sight `json:"sights"`
}

type SightComments struct {
	Sight Sight     `json:"sight"`
	Comms []Comment `json:"comments"`
}
