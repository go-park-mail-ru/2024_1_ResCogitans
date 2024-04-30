package storage

import (
	"github.com/go-park-mail-ru/2024_1_ResCogitans/internal/entities"
)

type SightStorageInterface interface {
	GetSightsList() ([]entities.Sight, error)
	GetSight(sightID int) (entities.Sight, error)
	SearchSights(str string) (entities.Sights, error)
	GetCommentsBySightID(commentID int) ([]entities.Comment, error)
	GetCommentsByUserID(userID int) ([]entities.Comment, error)
	CreateCommentBySightID(sightID int, comment entities.CommentRequest) error
	EditComment(commentID int, comment entities.CommentRequest) error
	DeleteComment(commentID int) error
	CreateJourney(journey entities.JourneyRequest) (entities.Journey, error)
	DeleteJourney(journeyID int) error
	GetJourneys(userID int) ([]entities.Journey, error)
	AddJourneySight(journeyID int, ids []int) error
	EditJourney(journeyID int, name, description string) error
	DeleteJourneySight(journeyID int, sight entities.SightIDRequest) error
	GetJourneySights(journeyID int) ([]entities.Sight, error)
	GetJourney(journeyID int) (entities.Journey, error)
}
