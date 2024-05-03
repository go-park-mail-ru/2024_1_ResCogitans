package entities

import (
	"time"
)

type ReviewRequest struct {
	Rating     int `json:"rating"`
	QuestionID int `json:"questionID"`
}

func (h ReviewRequest) Validate() error {
	return nil
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

type Review struct {
	ID         int
	UserID     int
	Rating     int
	QuestionID int
	CreatedAt  time.Time
}

type DataCheck struct {
	Questions []QuestionResponse `json:"questions"`
	Flag      bool               `json:"flag"`
}

type QuestionResponse struct {
	QuestionID int    `json:"questionID"`
	Text       string `json:"text"`
}
