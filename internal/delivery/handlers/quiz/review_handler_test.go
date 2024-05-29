package quiz_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/go-park-mail-ru/2024_1_ResCogitans/internal/delivery/handlers/quiz"
	"github.com/go-park-mail-ru/2024_1_ResCogitans/internal/entities"
	"github.com/go-park-mail-ru/2024_1_ResCogitans/internal/mocks"
	httperrors "github.com/go-park-mail-ru/2024_1_ResCogitans/utils/errors"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

func TestQuizHandler_CreateReview(t *testing.T) {
	mockQuestionUseCase := new(mocks.QuestionMockUseCase)
	mockCommentUseCase := new(mocks.CommentMockUseCase)
	mockJourneyUseCase := new(mocks.JourneyMockUseCase)

	handler := quiz.NewQuizHandler(mockQuestionUseCase, mockCommentUseCase, mockJourneyUseCase)

	t.Run("Error getting user id from context", func(t *testing.T) {
		response, err := handler.CreateReview(context.Background(), entities.Review{})

		assert.Error(t, err)
		assert.Equal(t, "failed getting user from context", err.Error())
		assert.Equal(t, false, response)
	})

	wrongCtx := context.WithValue(context.Background(), "userID", 0)

	t.Run("Unauthorized user", func(t *testing.T) {
		response, err := handler.CreateReview(wrongCtx, entities.Review{})

		assert.Error(t, err)
		httpError := httperrors.UnwrapHttpError(err)
		assert.Equal(t, "permission denied", httpError.Message)
		assert.Equal(t, http.StatusUnauthorized, httpError.Code)
		assert.Equal(t, false, response)
	})

	ctx := context.WithValue(context.Background(), "userID", 1)

	t.Run("Error creating review", func(t *testing.T) {
		mockQuestionUseCase.On("CreateReview", ctx, 1, entities.Review{UserID: 1, Rating: 4, QuestionID: 4}).Return(errors.New("error creating review")).Once()
		response, err := handler.CreateReview(ctx, entities.Review{UserID: 1, Rating: 4, QuestionID: 4})

		assert.Error(t, err)
		assert.Equal(t, "error creating review", err.Error())
		assert.Equal(t, false, response)

		mockQuestionUseCase.AssertExpectations(t)
	})

	mockQuestionUseCase.Mock.ExpectedCalls = nil

	t.Run("Successfully create review", func(t *testing.T) {
		mockQuestionUseCase.On("CreateReview", ctx, 1, entities.Review{UserID: 1, Rating: 4, QuestionID: 4}).Return(nil).Once()
		response, err := handler.CreateReview(ctx, entities.Review{UserID: 1, Rating: 4, QuestionID: 4})

		assert.NoError(t, err)
		assert.Equal(t, true, response)
		mockQuestionUseCase.AssertExpectations(t)
	})
}

func TestQuizHandler_CheckData(t *testing.T) {
	mockQuestionUseCase := new(mocks.QuestionMockUseCase)
	mockCommentUseCase := new(mocks.CommentMockUseCase)
	mockJourneyUseCase := new(mocks.JourneyMockUseCase)

	handler := quiz.NewQuizHandler(mockQuestionUseCase, mockCommentUseCase, mockJourneyUseCase)

	t.Run("Error getting user id from context", func(t *testing.T) {
		response, err := handler.CheckData(context.Background(), entities.Review{})

		assert.Error(t, err)
		assert.Equal(t, "failed getting user from context", err.Error())
		assert.Equal(t, entities.DataCheck{}, response)
	})

	wrongCtx := context.WithValue(context.Background(), "userID", 0)

	t.Run("Unauthorized user", func(t *testing.T) {
		response, err := handler.CheckData(wrongCtx, entities.Review{})

		assert.NoError(t, err)
		assert.Equal(t, entities.DataCheck{Flag: false}, response)
	})

	ctx := context.WithValue(context.Background(), "userID", 1)

	t.Run("Error getting question", func(t *testing.T) {
		mockQuestionUseCase.On("GetQuestions", ctx).Return([]entities.QuestionResponse{}, errors.New("error getting questions")).Once()
		response, err := handler.CheckData(ctx, entities.Review{})

		assert.Error(t, err)
		assert.Equal(t, "error getting questions", err.Error())
		assert.Equal(t, entities.DataCheck{}, response)
		mockQuestionUseCase.AssertExpectations(t)
	})

	mockQuestionUseCase.Mock.ExpectedCalls = nil

	t.Run("Review is already completed", func(t *testing.T) {
		mockQuestionUseCase.On("GetQuestions", ctx).Return([]entities.QuestionResponse{}, nil).Once()
		mockQuestionUseCase.On("CheckReview", ctx, 1).Return(true, nil).Once()
		response, err := handler.CheckData(ctx, entities.Review{})

		assert.NoError(t, err)
		assert.Equal(t, entities.DataCheck{Flag: false}, response)
		mockQuestionUseCase.AssertExpectations(t)
	})

	mockQuestionUseCase.Mock.ExpectedCalls = nil

	t.Run("Error checking review", func(t *testing.T) {
		mockQuestionUseCase.On("GetQuestions", ctx).Return([]entities.QuestionResponse{}, nil).Once()
		mockQuestionUseCase.On("CheckReview", ctx, 1).Return(false, errors.New("error checking review")).Once()
		response, err := handler.CheckData(ctx, entities.Review{})

		assert.Error(t, err)
		assert.Equal(t, "error checking review", err.Error())
		assert.Equal(t, entities.DataCheck{Flag: false}, response)
		mockQuestionUseCase.AssertExpectations(t)
	})

	mockQuestionUseCase.Mock.ExpectedCalls = nil

	t.Run("Journey was created", func(t *testing.T) {
		mockQuestionUseCase.On("GetQuestions", ctx).Return([]entities.QuestionResponse{}, nil).Once()
		mockQuestionUseCase.On("CheckReview", ctx, 1).Return(false, nil).Once()
		mockJourneyUseCase.On("CheckJourney", ctx, 1).Return(true, nil).Once()
		response, err := handler.CheckData(ctx, entities.Review{})

		assert.NoError(t, err)
		assert.Equal(t, entities.DataCheck{Questions: []entities.QuestionResponse{}, Flag: true}, response)

		mockQuestionUseCase.AssertExpectations(t)
		mockJourneyUseCase.AssertExpectations(t)
	})

	mockQuestionUseCase.Mock.ExpectedCalls = nil
	mockJourneyUseCase.Mock.ExpectedCalls = nil

	t.Run("Error checking journey", func(t *testing.T) {
		mockQuestionUseCase.On("GetQuestions", ctx).Return([]entities.QuestionResponse{}, nil).Once()
		mockQuestionUseCase.On("CheckReview", ctx, 1).Return(false, nil).Once()
		mockJourneyUseCase.On("CheckJourney", ctx, 1).Return(false, errors.New("error checking journey")).Once()
		response, err := handler.CheckData(ctx, entities.Review{})

		assert.Error(t, err)
		assert.Equal(t, "error checking journey", err.Error())
		assert.Equal(t, entities.DataCheck{Flag: false}, response)

		mockQuestionUseCase.AssertExpectations(t)
		mockJourneyUseCase.AssertExpectations(t)
	})

	mockQuestionUseCase.Mock.ExpectedCalls = nil
	mockJourneyUseCase.Mock.ExpectedCalls = nil

	t.Run("Comment was created", func(t *testing.T) {
		mockQuestionUseCase.On("GetQuestions", ctx).Return([]entities.QuestionResponse{}, nil).Once()
		mockQuestionUseCase.On("CheckReview", ctx, 1).Return(false, nil).Once()
		mockJourneyUseCase.On("CheckJourney", ctx, 1).Return(false, nil).Once()
		mockCommentUseCase.On("CheckCommentByUserID", ctx, 1).Return(true, nil).Once()
		response, err := handler.CheckData(ctx, entities.Review{})

		assert.NoError(t, err)
		assert.Equal(t, entities.DataCheck{Questions: []entities.QuestionResponse{}, Flag: true}, response)

		mockQuestionUseCase.AssertExpectations(t)
		mockJourneyUseCase.AssertExpectations(t)
		mockCommentUseCase.AssertExpectations(t)
	})

	mockQuestionUseCase.Mock.ExpectedCalls = nil
	mockJourneyUseCase.Mock.ExpectedCalls = nil
	mockCommentUseCase.Mock.ExpectedCalls = nil

	t.Run("Error checking comments", func(t *testing.T) {
		mockQuestionUseCase.On("GetQuestions", ctx).Return([]entities.QuestionResponse{}, nil).Once()
		mockQuestionUseCase.On("CheckReview", ctx, 1).Return(false, nil).Once()
		mockJourneyUseCase.On("CheckJourney", ctx, 1).Return(false, nil).Once()
		mockCommentUseCase.On("CheckCommentByUserID", ctx, 1).Return(false, errors.New("error checking comments")).Once()
		response, err := handler.CheckData(ctx, entities.Review{})

		assert.Error(t, err)
		assert.Equal(t, "error checking comments", err.Error())
		assert.Equal(t, entities.DataCheck{Flag: false}, response)

		mockQuestionUseCase.AssertExpectations(t)
		mockJourneyUseCase.AssertExpectations(t)
		mockCommentUseCase.AssertExpectations(t)
	})

	mockQuestionUseCase.Mock.ExpectedCalls = nil
	mockJourneyUseCase.Mock.ExpectedCalls = nil
	mockCommentUseCase.Mock.ExpectedCalls = nil

	t.Run("Not active user", func(t *testing.T) {
		mockQuestionUseCase.On("GetQuestions", ctx).Return([]entities.QuestionResponse{}, nil).Once()
		mockQuestionUseCase.On("CheckReview", ctx, 1).Return(false, nil).Once()
		mockJourneyUseCase.On("CheckJourney", ctx, 1).Return(false, nil).Once()
		mockCommentUseCase.On("CheckCommentByUserID", ctx, 1).Return(false, nil).Once()
		response, err := handler.CheckData(ctx, entities.Review{})

		assert.NoError(t, err)
		assert.Equal(t, entities.DataCheck{Flag: false}, response)

		mockQuestionUseCase.AssertExpectations(t)
		mockJourneyUseCase.AssertExpectations(t)
		mockCommentUseCase.AssertExpectations(t)
	})

	mockQuestionUseCase.Mock.ExpectedCalls = nil
	mockJourneyUseCase.Mock.ExpectedCalls = nil
	mockCommentUseCase.Mock.ExpectedCalls = nil
}

func TestQuizHandler_SetStat(t *testing.T) {
	mockQuestionUseCase := new(mocks.QuestionMockUseCase)
	mockCommentUseCase := new(mocks.CommentMockUseCase)
	mockJourneyUseCase := new(mocks.JourneyMockUseCase)

	handler := quiz.NewQuizHandler(mockQuestionUseCase, mockCommentUseCase, mockJourneyUseCase)

	t.Run("Error getting request from context", func(t *testing.T) {
		response, err := handler.SetStat(context.Background(), entities.Statistic{})

		assert.Error(t, err)
		assert.Equal(t, "failed getting user from context", err.Error())
		assert.Equal(t, []entities.Statistic{}, response)
	})

	wrongCtx := context.WithValue(context.Background(), "userID", 0)

	t.Run("Unauthorized user", func(t *testing.T) {
		response, err := handler.SetStat(wrongCtx, entities.Statistic{})

		assert.Error(t, err)
		httpError := httperrors.UnwrapHttpError(err)
		assert.Equal(t, "permission denied", httpError.Message)
		assert.Equal(t, http.StatusUnauthorized, httpError.Code)
		assert.Equal(t, []entities.Statistic{}, response)
	})

	ctx := context.WithValue(context.Background(), "userID", 1)

	t.Run("Error setting stat", func(t *testing.T) {
		mockQuestionUseCase.On("SetStat", ctx, 1).Return([]entities.Statistic{}, errors.New("error setting stat")).Once()
		response, err := handler.SetStat(ctx, entities.Statistic{})

		assert.Error(t, err)
		assert.Equal(t, "error setting stat", err.Error())
		assert.Equal(t, []entities.Statistic{}, response)

		mockQuestionUseCase.AssertExpectations(t)
	})

	mockQuestionUseCase.Mock.ExpectedCalls = nil

	t.Run("Success setting stat", func(t *testing.T) {
		mockQuestionUseCase.On("SetStat", ctx, 1).Return([]entities.Statistic{{}}, nil)
		response, err := handler.SetStat(ctx, entities.Statistic{})

		assert.NoError(t, err)
		assert.Equal(t, []entities.Statistic{{}}, response)
		mockQuestionUseCase.AssertExpectations(t)
	})
}
