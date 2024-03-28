package storage

import (
	"errors"
	"sync"
)

type StorageInterface interface {
	SaveSession(sessionID string, userID int)
	GetSession(sessionID string) (int, error)
	DeleteSession(sessionID string)
}

type SessionStorage struct {
	Store map[string]int
	mu    sync.Mutex
}

func NewSessionStorage() StorageInterface {
	return &SessionStorage{
		Store: make(map[string]int),
	}
}

func (ss *SessionStorage) SaveSession(sessionID string, userID int) {
	ss.mu.Lock()
	defer ss.mu.Unlock()
	ss.Store[sessionID] = userID
}

func (ss *SessionStorage) GetSession(sessionID string) (int, error) {
	ss.mu.Lock()
	defer ss.mu.Unlock()
	userID, ok := ss.Store[sessionID]
	if !ok {
		return 0, errors.New("Session not found")
	}
	return userID, nil
}

func (ss *SessionStorage) DeleteSession(sessionID string) {
	ss.mu.Lock()
	defer ss.mu.Unlock()
	delete(ss.Store, sessionID)
}
