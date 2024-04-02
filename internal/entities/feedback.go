package entities

type Comment struct {
	ID       int    `json:"id"`
	UserID   int    `json:"userID"`
	Username string `json:"username"`
	SightID  int    `json:"sightID"`
	Rating   int    `json:"rating"`
	Feedback string `json:"feedback"`
}

func (h Comment) Validate() error {
	return nil
}
