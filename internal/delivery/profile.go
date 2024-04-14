package delivery

import (
	"context"
	"encoding/json"
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

	errProfileResetPassword = errors.HttpError{
		Code:    http.StatusUnauthorized,
		Message: "weak password",
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
		logger.Error("Error while converting string to int", "error", err)
		return entities.UserProfile{}, errParsing
	}
	r, _ := httputils.HttpRequest(ctx)
	if userID != usecase.GetSession(r) {
		logger.Error("Cannot edit other's profile")
		return entities.UserProfile{}, errProfilePermissionDenied
	}

	dataInt := make(map[string]int)
	dataStr := make(map[string]string)

	dataInt["userID"] = userID
	dataStr["username"] = requestData.Username
	dataStr["bio"] = requestData.Bio

	userRepo := userRep.NewUserRepo(db)
	profile, err := userRepo.EditUserProfile(dataInt, dataStr)

	if err != nil {
		return entities.UserProfile{}, errEditProfile
	}

	return profile, nil
}

func (h *ProfileHandler) UpdateUserPassword(ctx context.Context, requestData entities.User) (ProfileResponse, error) {
	logger := logger.Logger()

	db, err := db.GetPostgres()
	if err != nil {
		logger.Error("Error while connecting to db", "error", err)
	}

	pathParams := wrapper.GetPathParamsFromCtx(ctx)
	userID, err := strconv.Atoi(pathParams["id"])
	if err != nil {
		logger.Error("Error while converting string to int", "error", err)
		return ProfileResponse{}, errParsing
	}
	r, _ := httputils.HttpRequest(ctx)
	if userID != usecase.GetSession(r) {
		logger.Error("Cannot edit other's profile")
		return ProfileResponse{}, errProfilePermissionDenied
	}

	if !entities.ValidatePassword(requestData.Passwrd) {
		return ProfileResponse{}, errProfileResetPassword
	}

	userRepo := userRep.NewUserRepo(db)
	err = userRepo.UpdateUserPassword(requestData.ID, requestData.Passwrd)

	if err != nil {
		return ProfileResponse{}, err
	}

	return ProfileResponse{}, nil
}

// TODO: нужно будет убрать это инженерное решение (костыль) после фикса обертки
func (f *ProfileHandler) UploadFile(w http.ResponseWriter, r *http.Request) {
	logger := logger.Logger()
	db, err := db.GetPostgres()
	if err != nil {
		logger.Error(err.Error())
		errors.WriteHttpError(err, w)
		return
	}

	userID, err := strconv.Atoi(wrapper.GetPathParams(r)["id"])
	if err != nil {
		logger.Error("Handler error", "error", err)
		errors.WriteHttpError(err, w)
		return
	}

	if userID != usecase.GetSession(r) {
		logger.Error("Cannot edit other's profile")
		errors.WriteHttpError(errProfilePermissionDenied, w)
		return
	}

	path, err := SaveFile(r)
	if err != nil {
		logger.Error("Handler error", "error", err)
		errors.WriteHttpError(err, w)
		return
	}

	dataInt := map[string]int{"userID": userID}
	dataStr := map[string]string{"avatar": path}

	userRepo := userRep.NewUserRepo(db)
	profile, err := userRepo.EditUserProfile(dataInt, dataStr)
	if err != nil {
		logger.Error("Handler error", "error", err)
		errors.WriteHttpError(err, w)
		return
	}

	rawJSON, err := json.Marshal(profile)
	if err != nil {
		logger.Error("Error encoding response", "error", err)
		errors.WriteHttpError(err, w)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(rawJSON)
}
