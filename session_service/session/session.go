package session

import (
	"context"

	"github.com/go-park-mail-ru/2024_1_ResCogitans/internal/service/gen"
)

type SessionService struct {
	// Поля, необходимые для сервиса сессий
}

func NewSessionService() *SessionService {
	return &SessionService{
		// Инициализация полей
	}
}

func (s *SessionService) CreateSession(ctx context.Context, req *gen.SaveSessionRequest) (*gen.SaveSessionResponse, error) {
	// Логика создания сессии
	// Возвращает SaveSessionResponse и ошибку
}

func (s *SessionService) GetSession(ctx context.Context, req *gen.GetSessionRequest) (*gen.GetSessionResponse, error) {
	// Логика получения сессии
	// Возвращает GetSessionResponse и ошибку
}

func (s *SessionService) DeleteSession(ctx context.Context, req *gen.DeleteSessionRequest) (*gen.DeleteSessionResponse, error) {
	// Логика удаления сессии
	// Возвращает DeleteSessionResponse и ошибку
}
