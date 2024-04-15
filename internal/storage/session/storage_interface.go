package session

type StorageInterface interface {
	SaveSession(sessionID string, userID int) error
	GetSession(sessionID string) (int, error)
	DeleteSession(sessionID string) error
}
