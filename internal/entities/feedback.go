package entities

type Comment struct {
	ID       int    `json:"id"`
	UserID   int    `json:"userID"`
	Email    string `json:"username"`
	SightID  int    `json:"sightID"`
	Rating   int    `json:"rating"`
	Feedback string `json:"feedback"`
}

type Comments struct {
	Comment []Comment `json:"comments"`
}

func (h Comment) Validate() error {
	return nil
}

func (h Comments) Validate() error {
	return nil
}
