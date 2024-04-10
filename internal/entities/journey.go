package entities

type Journey struct {
	ID          int    `json:"id"`
	UserID      int    `json:"userID"`
	Email       string `json:"username"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type JourneySight struct {
	ID        int    `json:"id"`
	UserID    int    `json:"userID"`
	JourneyID int    `json:"journeyID"`
	SightID   int    `json:"sightID"`
	Priority  int    `json:"priority"`
	SightName string `json:"sight_name"`
}

type Journeys struct {
	Journey []Journey `json:"journeys"`
}

type JourneySights struct {
	Journey Journey `json:"journey"`
	Sight   []Sight `json:"sights"`
}

func (h Journey) Validate() error {
	return nil
}

func (h JourneySight) Validate() error {
	return nil
}
