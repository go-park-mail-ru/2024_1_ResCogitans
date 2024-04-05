package delivery

import (
	"context"
	"errors"
	"strconv"

	"github.com/go-park-mail-ru/2024_1_ResCogitans/internal/entities"
	"github.com/go-park-mail-ru/2024_1_ResCogitans/internal/http-server/server/db"
	sightRep "github.com/go-park-mail-ru/2024_1_ResCogitans/internal/repository/postgres"
	"github.com/go-park-mail-ru/2024_1_ResCogitans/utils/logger"
	"github.com/go-park-mail-ru/2024_1_ResCogitans/utils/wrapper"
)

type CommentHandler struct{}

func (h *CommentHandler) CreateComment(ctx context.Context, requestData entities.Comment) (entities.Comment, error) {
	db, err := db.GetPostgres()

	if err != nil {
		logger.Logger().Error(err.Error())
	}

	pathParams := wrapper.GetPathParamsFromCtx(ctx)
	sightID, err := strconv.Atoi(pathParams["id"])
	if err != nil {
		logger.Logger().Error("Cannot convert string to integer to get sight")
		return entities.Comment{}, err
	}

	dataStr := make(map[string]string)
	dataInt := make(map[string]int)

	dataInt["userID"] = requestData.UserID
	dataInt["sightID"] = sightID
	dataInt["rating"] = requestData.Rating

	dataStr["feedback"] = requestData.Feedback

	sightsRepo := sightRep.NewSightRepo(db)
	err = sightsRepo.CreateCommentBySightID(dataStr, dataInt)

	if err != nil {
		return entities.Comment{}, errors.New("cannot create comment")
	}

	return entities.Comment{}, nil
}

func (h *CommentHandler) EditComment(ctx context.Context, requestData entities.Comment) (entities.Comment, error) {
	db, err := db.GetPostgres()

	if err != nil {
		logger.Logger().Error(err.Error())
	}

	pathParams := wrapper.GetPathParamsFromCtx(ctx)
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

	sightsRepo := sightRep.NewSightRepo(db)
	err = sightsRepo.EditCommentByCommentID(dataStr, dataInt)

	if err != nil {
		return entities.Comment{}, errors.New("cannot edit comment")
	}

	return entities.Comment{}, nil
}

func (h *CommentHandler) DeleteComment(ctx context.Context, requestData entities.Comment) (entities.Comment, error) {
	db, err := db.GetPostgres()

	if err != nil {
		logger.Logger().Error(err.Error())
	}

	pathParams := wrapper.GetPathParamsFromCtx(ctx)
	commentID, err := strconv.Atoi(pathParams["cid"])
	if err != nil {
		logger.Logger().Error("Cannot convert string to integer to get sight")
		return entities.Comment{}, err
	}

	dataInt := make(map[string]int)
	dataInt["id"] = commentID

	sightsRepo := sightRep.NewSightRepo(db)
	err = sightsRepo.DeleteCommentByCommentID(dataInt)

	if err != nil {
		return entities.Comment{}, errors.New("cannot edit comment")
	}

	return entities.Comment{}, nil
}
