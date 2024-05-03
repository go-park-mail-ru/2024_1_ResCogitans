package quiz

import (
	"context"
	"net/http"

	"github.com/go-park-mail-ru/2024_1_ResCogitans/internal/entities"
	"github.com/go-park-mail-ru/2024_1_ResCogitans/internal/usecase"
	httperrors "github.com/go-park-mail-ru/2024_1_ResCogitans/utils/errors"
	"github.com/go-park-mail-ru/2024_1_ResCogitans/utils/httputils"
)

type QuizHandler struct {
	questionUseCase usecase.QuestionUseCaseInterface
	commentUseCase  usecase.CommentUseCaseInterface
	journeyUseCase  usecase.JourneyUseCaseInterface
}

func NewQuizHandler(questionUseCase usecase.QuestionUseCaseInterface,
	commentUseCase usecase.CommentUseCaseInterface,
	journeyUseCase usecase.JourneyUseCaseInterface) *QuizHandler {
	return &QuizHandler{
		questionUseCase: questionUseCase,
		commentUseCase:  commentUseCase,
		journeyUseCase:  journeyUseCase,
	}
}

// CreateReview godoc
// @Summary Создание отзыва сайта
// @Description Возвращает true при удачном создании, false при неудачном
// @Tags Опросник
// @Accept json
// @Produce json
// @Param user body entities.ReviewRequest true "Отзыв пользователя и номер вопроса"
// @Success 200 {object} bool
// @Failure 401 {object} httperrors.HttpError
// @Failure 500 {object} httperrors.HttpError
// @Router /api/review/create [post]
func (h *QuizHandler) CreateReview(ctx context.Context, requestData entities.ReviewRequest) (bool, error) {
	userID, err := httputils.GetUserFromCtx(ctx)
	if err != nil {
		return false, err
	}
	if userID == 0 {
		return false, httperrors.NewHttpError(http.StatusUnauthorized, "Permission denied")
	}

	err = h.questionUseCase.CreateReview(userID, requestData)
	if err != nil {
		return false, err
	}
	return true, nil
}

// CheckData godoc
// @Summary Проверка необходимости отображения опросника
// @Tags Опросник
// @Accept json
// @Produce json
// @Success 200 {object} entities.DataCheck
// @Failure 500 {object} httperrors.HttpError
// @Router /api/review/check [get]
func (h *QuizHandler) CheckData(ctx context.Context, _ entities.ReviewRequest) (entities.DataCheck, error) {
	userID, err := httputils.GetUserFromCtx(ctx)
	if err != nil {
		return entities.DataCheck{}, err
	}

	if userID == 0 {
		return entities.DataCheck{Flag: false}, nil
	}

	questions, err := h.questionUseCase.GetQuestions()
	if err != nil {
		return entities.DataCheck{}, err
	}

	ok, err := h.questionUseCase.CheckReview(userID)
	if err != nil {
		return entities.DataCheck{}, err
	}
	if ok {
		return entities.DataCheck{Flag: false}, nil
	}

	ok, err = h.journeyUseCase.CheckJourney(userID)
	if err != nil {
		return entities.DataCheck{}, err
	}
	if ok {
		return entities.DataCheck{Flag: true, Questions: questions}, nil
	}

	ok, err = h.commentUseCase.CheckCommentByUserID(userID)
	if err != nil {
		return entities.DataCheck{}, err
	}
	if ok {
		return entities.DataCheck{Flag: true, Questions: questions}, nil
	}
	return entities.DataCheck{Flag: false}, nil
}

// SetStat godoc
// @Summary Статистика опросника для конкретного пользователя
// @Tags Опросник
// @Accept json
// @Produce json
// @Success 200 {object} []entities.Statistic
// @Failure 401 {object} httperrors.HttpError
// @Failure 500 {object} httperrors.HttpError
// @Router /api/review/get [get]
func (h *QuizHandler) SetStat(ctx context.Context, _ entities.Statistic) ([]entities.Statistic, error) {
	userID, err := httputils.GetUserFromCtx(ctx)
	if err != nil {
		return []entities.Statistic{}, err
	}
	if userID == 0 {
		return []entities.Statistic{}, httperrors.NewHttpError(http.StatusUnauthorized, "Permission denied")
	}

	stat, err := h.questionUseCase.SetStat(userID)
	if err != nil {
		return []entities.Statistic{}, err
	}

	return stat, nil
}
