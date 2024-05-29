package mocks

import (
	"context"

	"github.com/go-park-mail-ru/2024_1_ResCogitans/internal/entities"
	"github.com/stretchr/testify/mock"
)

type CommentMockUseCase struct {
	mock.Mock
}

func (m *CommentMockUseCase) CreateCommentBySightID(ctx context.Context, sightID int, comment entities.Comment) error {
	args := m.Called(ctx, sightID, comment)
	return args.Error(0)
}

func (m *CommentMockUseCase) EditCommentByCommentID(ctx context.Context, commentID int, comment entities.Comment) error {
	args := m.Called(ctx, commentID, comment)
	return args.Error(0)
}

func (m *CommentMockUseCase) DeleteCommentByCommentID(ctx context.Context, commentID int) error {
	args := m.Called(ctx, commentID)
	return args.Error(0)
}

func (m *CommentMockUseCase) CheckCommentByUserID(ctx context.Context, userID int) (bool, error) {
	args := m.Called(ctx, userID)
	return args.Bool(0), args.Error(1)
}
