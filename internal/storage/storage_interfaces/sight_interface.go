package storage

import (
	"github.com/go-park-mail-ru/2024_1_ResCogitans/internal/entities"
)

type SightStorageInterface interface {
	GetSightsList() ([]entities.Sight, error)
	GetSight(sightID int) (entities.Sight, error)
	GetCommentsBySightID(commentID int) ([]entities.Comment, error)
	CreateCommentBySightID(dataStr map[string]string, dataInt map[string]int) error
	EditComment(dataStr map[string]string, dataInt map[string]int) error
	DeleteComment(dataInt map[string]int) error
	CreateJourney(dataInt map[string]int, dataStr map[string]string) (entities.Journey, error)
	DeleteJourney(dataInt map[string]int) error
	GetJourneys(userID int) ([]entities.Journey, error)
	AddJourneySight(journeyID int, ids []int) error
	DeleteJourneySight(dataInt map[string]int) error
	GetJourneySights(journeyID int) ([]entities.Sight, error)
	GetJourney(journeyID int) (entities.Journey, error)
}
