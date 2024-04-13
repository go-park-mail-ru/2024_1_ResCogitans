package delivery

import (
	"context"
	"net/http"
	"strconv"

	"github.com/go-park-mail-ru/2024_1_ResCogitans/internal/entities"
	"github.com/go-park-mail-ru/2024_1_ResCogitans/internal/http-server/server/db"
	userRep "github.com/go-park-mail-ru/2024_1_ResCogitans/internal/repository/postgres"
	"github.com/go-park-mail-ru/2024_1_ResCogitans/internal/usecase"
	"github.com/go-park-mail-ru/2024_1_ResCogitans/utils/errors"
	"github.com/go-park-mail-ru/2024_1_ResCogitans/utils/httputils"
	"github.com/go-park-mail-ru/2024_1_ResCogitans/utils/logger"
	"github.com/go-park-mail-ru/2024_1_ResCogitans/utils/wrapper"
)

type ProfileHandler struct{}

type ProfileResponse struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Bio      string `json:"bio"`
	Avatar   string `json:"avatar"`
}

var (
	errEditProfile = errors.HttpError{
		Code:    http.StatusInternalServerError,
		Message: "failed edit profile",
	}
	errDeleteProfile = errors.HttpError{
		Code:    http.StatusInternalServerError,
		Message: "failed deleting profile",
	}
	errProfilePermissionDenied = errors.HttpError{
		Code:    http.StatusUnauthorized,
		Message: "permission denied",
	}
)

func (h *ProfileHandler) GetUserProfile(ctx context.Context, requestData entities.User) (ProfileResponse, error) {
	db, err := db.GetPostgres()
	if err != nil {
		logger.Logger().Error(err.Error())
	}

	pathParams := wrapper.GetPathParamsFromCtx(ctx)
	id, err := strconv.Atoi(pathParams["id"])
	if err != nil {
		logger.Logger().Error("Cannot convert string to integer to get sight")
		return ProfileResponse{}, err
	}

	dataInt := make(map[string]int)

	dataInt["userID"] = id

	UserRepo := userRep.NewUserRepo(db)
	user, err := UserRepo.GetUserProfile(dataInt)
	if err != nil {
		return ProfileResponse{}, errLoginUser
	}

	profileResponse := ProfileResponse{
		ID:       user.UserID,
		Username: user.Username,
		Bio:      user.Bio,
		Avatar:   user.Avatar,
	}

	return profileResponse, nil
}

func (h *ProfileHandler) DeleteUser(ctx context.Context, requestData entities.User) (ProfileResponse, error) {
	logger := logger.Logger()
	db, err := db.GetPostgres()

	if err != nil {
		logger.Error(err.Error())
	}

	pathParams := wrapper.GetPathParamsFromCtx(ctx)
	userID, err := strconv.Atoi(pathParams["id"])
	if err != nil {
		logger.Error("Cannot convert string to integer to get sight")
		return ProfileResponse{}, errParsing
	}

	if r, _ := httputils.HttpRequest(ctx); userID != usecase.GetSession(r) {
		logger.Error("Cannot edit other's profile")
		return ProfileResponse{}, errProfilePermissionDenied
	}

	dataInt := make(map[string]int)

	dataInt["userID"] = userID

	userRepo := userRep.NewUserRepo(db)
	err = userRepo.DeleteUserProfile(dataInt)

	if err != nil {
		return ProfileResponse{}, errDeleteProfile
	}

	return ProfileResponse{}, nil
}

func (h *ProfileHandler) EditUserProfile(ctx context.Context, requestData entities.UserProfile) (entities.UserProfile, error) {
	logger := logger.Logger()
	db, err := db.GetPostgres()

	if err != nil {
		logger.Error(err.Error())
	}

	pathParams := wrapper.GetPathParamsFromCtx(ctx)
	userID, err := strconv.Atoi(pathParams["id"])
	if err != nil {
		logger.Error("Cannot convert string to integer to get sight")
		return entities.UserProfile{}, errParsing
	}

	if r, _ := httputils.HttpRequest(ctx); userID != usecase.GetSession(r) {
		logger.Error("Cannot edit other's profile")
		return entities.UserProfile{}, errProfilePermissionDenied
	}

	dataInt := make(map[string]int)
	dataStr := make(map[string]string)

	dataInt["userID"] = userID
	dataStr["username"] = requestData.Username
	dataStr["bio"] = requestData.Bio
	dataStr["avatar"] = requestData.Avatar

	userRepo := userRep.NewUserRepo(db)
	profile, err := userRepo.EditUserProfile(dataInt, dataStr)

	if err != nil {
		return entities.UserProfile{}, errEditProfile
	}

	return profile, nil
}
