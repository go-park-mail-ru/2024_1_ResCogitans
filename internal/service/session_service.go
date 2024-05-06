package service

import (
	storage "github.com/go-park-mail-ru/2024_1_ResCogitans/internal/storage/storage_interfaces"
)

type SessionService struct {
	storage storage.SessionStorageInterface
}

func NewSessionService(storage storage.SessionStorageInterface) *SessionService {
	return &SessionService{storage: storage}
}

func (s *SessionService) SaveSession(sessionID string, userID int) error {
	return s.storage.SaveSession(sessionID, userID)
}

func (s *SessionService) GetSession(sessionID string) (int, error) {
	return s.storage.GetSession(sessionID)
}

func (s *SessionService) DeleteSession(sessionID string) error {
	return s.storage.DeleteSession(sessionID)
}
