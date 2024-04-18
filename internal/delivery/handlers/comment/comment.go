package comment

import (
	"context"
	"strconv"

	"github.com/go-park-mail-ru/2024_1_ResCogitans/internal/entities"
	"github.com/go-park-mail-ru/2024_1_ResCogitans/internal/usecase"
	"github.com/go-park-mail-ru/2024_1_ResCogitans/utils/httputils"
	"github.com/go-park-mail-ru/2024_1_ResCogitans/utils/logger"
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

	dataStr := make(map[string]string)
	dataInt := make(map[string]int)

	dataInt["userID"] = requestData.UserID
	dataInt["sightID"] = sightID
	dataInt["rating"] = requestData.Rating

	dataStr["feedback"] = requestData.Feedback

	err = h.CommentUseCase.CreateCommentBySightID(dataStr, dataInt)
	if err != nil {
		return entities.Comment{}, err
	}

	return entities.Comment{}, nil
}

func (h *CommentHandler) EditComment(ctx context.Context, requestData entities.Comment) (entities.Comment, error) {
	pathParams := httputils.GetPathParamsFromCtx(ctx)
	commentID, err := strconv.Atoi(pathParams["cid"])
	if err != nil {
		logger.Logger().Error("Cannot convert string to integer to get sight")
		return entities.Comment{}, err
	}

	dataStr := make(map[string]string)
	dataInt := make(map[string]int)

	dataInt["id"] = commentID
	dataInt["rating"] = requestData.Rating
	dataStr["feedback"] = requestData.Feedback

	err = h.CommentUseCase.EditCommentByCommentID(dataStr, dataInt)
	if err != nil {
		return entities.Comment{}, err
	}

	return entities.Comment{}, nil
}

func (h *CommentHandler) DeleteComment(ctx context.Context, _ entities.Comment) (entities.Comment, error) {
	pathParams := httputils.GetPathParamsFromCtx(ctx)
	commentID, err := strconv.Atoi(pathParams["cid"])
	if err != nil {
		logger.Logger().Error("Cannot convert string to integer to get sight")
		return entities.Comment{}, err
	}

	dataInt := make(map[string]int)
	dataInt["id"] = commentID

	err = h.CommentUseCase.DeleteCommentByCommentID(dataInt)
	if err != nil {
		return entities.Comment{}, err
	}

	return entities.Comment{}, nil
}
