package entities

import (
	"time"
)

type Review struct {
	ID         int       `json:"id"`
	UserID     int       `json:"userID"`
	Rating     int       `json:"rating"`
	QuestionID int       `json:"questionID"`
	CreatedAt  time.Time `json:"createdAt"`
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

type Statistic struct {
	ID           int     `json:"id"`
	Text         string  `json:"text"`
	UserGrade    int     `json:"UserGrade"`
	AverageGrade float64 `json:"AverageGrade"`
}

func (h Statistic) Validate() error {
	return nil
}
