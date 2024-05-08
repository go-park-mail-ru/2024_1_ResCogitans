package comment

import (
	"context"
	"strconv"

	"github.com/go-park-mail-ru/2024_1_ResCogitans/internal/entities"
	"github.com/go-park-mail-ru/2024_1_ResCogitans/internal/usecase"
	"github.com/go-park-mail-ru/2024_1_ResCogitans/utils/httputils"
	"github.com/pkg/errors"
)

type CommentHandler struct {
	CommentUseCase usecase.CommentUseCaseInterface
}

func NewCommentHandler(usecase usecase.CommentUseCaseInterface) *CommentHandler {
	return &CommentHandler{
		CommentUseCase: usecase,
	}
}

func (h *CommentHandler) CreateComment(ctx context.Context, requestData entities.Comment) (entities.Comment, error) {
	pathParams := httputils.GetPathParamsFromCtx(ctx)
	sightID, err := strconv.Atoi(pathParams["id"])
	if err != nil {
		return entities.Comment{}, err
	}

	if requestData.Rating < 0 {
		return entities.Comment{}, errors.New("Rating must be greater than zero")
	}

	err = h.CommentUseCase.CreateCommentBySightID(sightID, requestData)
	if err != nil {
		return entities.Comment{}, err
	}

	return entities.Comment{}, nil
}

func (h *CommentHandler) EditComment(ctx context.Context, requestData entities.Comment) (entities.Comment, error) {
	pathParams := httputils.GetPathParamsFromCtx(ctx)
	commentID, err := strconv.Atoi(pathParams["cid"])
	if err != nil {
		return entities.Comment{}, err
	}

	err = h.CommentUseCase.EditCommentByCommentID(commentID, requestData)
	if err != nil {
		return entities.Comment{}, err
	}

	return entities.Comment{}, nil
}

func (h *CommentHandler) DeleteComment(ctx context.Context, _ entities.Comment) (entities.Comment, error) {
	pathParams := httputils.GetPathParamsFromCtx(ctx)
	commentID, err := strconv.Atoi(pathParams["cid"])
	if err != nil {
		return entities.Comment{}, err
	}

	err = h.CommentUseCase.DeleteCommentByCommentID(commentID)
	if err != nil {
		return entities.Comment{}, err
	}

	return entities.Comment{}, nil
}
