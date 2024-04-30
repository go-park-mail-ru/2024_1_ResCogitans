package entities

type CommentRequest struct {
	UserID   int    `json:"userID"`
	Rating   int    `json:"rating"`
	Feedback string `json:"feedback"`
}

func (h CommentRequest) Validate() error {
	return nil
}

type Comment struct {
	ID       int
	UserID   int
	Username string
	SightID  int
	Rating   int
	Feedback string
	Avatar   string
}

type Comments struct {
	Comment []Comment `json:"comments"`
}

func (h Comments) Validate() error {
	return nil
}
