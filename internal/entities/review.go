package entities

type Review struct {
	ID         int    `json:"id"`
	UserID     int    `json:"userID"`
	Rating     int    `json:"rating"`
	QuestionID int    `json:"questionID"`
	CreatedAt  string `json:"createdAt"`
}

func (h Review) Validate() error {
	return nil
}

type DataCheck struct {
	Questions []QuestionResponse `json:"questions"`
	Flag      bool               `json:"flag"`
}

type QuestionResponse struct {
	QuestionID int    `json:"questionID"`
	Text       string `json:"text"`
}
