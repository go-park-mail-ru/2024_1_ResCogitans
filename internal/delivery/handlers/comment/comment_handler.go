package comment

import (
	"context"
	"net/http"
	"strconv"

	"github.com/go-park-mail-ru/2024_1_ResCogitans/internal/entities"
	"github.com/go-park-mail-ru/2024_1_ResCogitans/internal/usecase"
	httperrors "github.com/go-park-mail-ru/2024_1_ResCogitans/utils/errors"
	"github.com/go-park-mail-ru/2024_1_ResCogitans/utils/httputils"
)

type CommentHandler struct {
	CommentUseCase usecase.CommentUseCaseInterface
}

func NewCommentHandler(usecase usecase.CommentUseCaseInterface) *CommentHandler {
	return &CommentHandler{
		CommentUseCase: usecase,
	}
}

// CreateComment godoc
// @Summary Создание отзыва
// @Tags Отзывы
// @Accept json
// @Produce json
// @Param user body entities.CommentRequest true "Данные отзыва"
// @Success 200 {object} bool
// @Failure 400 {object} httperrors.HttpError
// @Failure 500 {object} httperrors.HttpError
// @Router /api/sight/{id}/create [post]
func (h *CommentHandler) CreateComment(ctx context.Context, requestData entities.CommentRequest) (bool, error) {
	pathParams := httputils.GetPathParamsFromCtx(ctx)
	sightID, err := strconv.Atoi(pathParams["id"])
	if err != nil {
		return false, err
	}

	if requestData.Rating < 0 {
		return false, httperrors.NewHttpError(http.StatusBadRequest, "Rating must be greater than zero")
	}

	err = h.CommentUseCase.CreateCommentBySightID(sightID, requestData)
	if err != nil {
		return false, err
	}
	return true, nil
}

// EditComment godoc
// @Summary Редактирование отзыва
// @Tags Отзывы
// @Accept json
// @Produce json
// @Param user body entities.CommentRequest true "Данные отзыва"
// @Success 200 {object} bool
// @Failure 400 {object} httperrors.HttpError
// @Failure 500 {object} httperrors.HttpError
// @Router /api/sight/{sid}/edit/{cid} [post]
func (h *CommentHandler) EditComment(ctx context.Context, requestData entities.CommentRequest) (bool, error) {
	pathParams := httputils.GetPathParamsFromCtx(ctx)
	commentID, err := strconv.Atoi(pathParams["cid"])
	if err != nil {
		return false, err
	}

	err = h.CommentUseCase.EditCommentByCommentID(commentID, requestData)
	if err != nil {
		return false, err
	}
	return true, nil
}

// DeleteComment godoc
// @Summary Удаление отзыва
// @Tags Отзывы
// @Accept json
// @Produce json
// @Success 200 {object} bool
// @Failure 400 {object} httperrors.HttpError
// @Failure 500 {object} httperrors.HttpError
// @Router /api/sight/{sid}/delete/{cid} [post]
func (h *CommentHandler) DeleteComment(ctx context.Context, _ entities.CommentRequest) (bool, error) {
	pathParams := httputils.GetPathParamsFromCtx(ctx)
	commentID, err := strconv.Atoi(pathParams["cid"])
	if err != nil {
		return false, err
	}

	err = h.CommentUseCase.DeleteCommentByCommentID(commentID)
	if err != nil {
		return false, err
	}

	return true, nil
}
