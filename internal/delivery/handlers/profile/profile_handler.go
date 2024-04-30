package profile

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-park-mail-ru/2024_1_ResCogitans/internal/delivery"
	"github.com/go-park-mail-ru/2024_1_ResCogitans/internal/entities"
	"github.com/go-park-mail-ru/2024_1_ResCogitans/internal/usecase"
	"github.com/go-park-mail-ru/2024_1_ResCogitans/utils/httputils"
	"github.com/go-park-mail-ru/2024_1_ResCogitans/utils/logger"
	"github.com/go-park-mail-ru/2024_1_ResCogitans/utils/wrapper"
)

type ProfileHandler struct {
	userProfileUseCase usecase.ProfileUseCaseInterface
}

func NewProfileHandler(userProfileUseCase usecase.ProfileUseCaseInterface) *ProfileHandler {
	return &ProfileHandler{
		userProfileUseCase: userProfileUseCase,
	}
}

// Get godoc
// @Summary Получение данных пользователя
// @Tags Профиль
// @Accept json
// @Produce json
// @Success 200 {object} entities.UserProfile
// @Failure 500 {object} httperrors.HttpError
// @Router /api/profile/{id} [get]
func (h *ProfileHandler) Get(ctx context.Context, _ entities.ProfileRequest) (entities.UserProfile, error) {
	pathParams := httputils.GetPathParamsFromCtx(ctx)
	idStr := pathParams["id"]
	userID, err := strconv.Atoi(idStr)
	if err != nil {
		return entities.UserProfile{}, err
	}
	return h.userProfileUseCase.GetUserProfile(userID)
}

// Edit godoc
// @Summary Изменение данных пользователя
// @Tags Профиль
// @Accept json
// @Param user body entities.ProfileRequest true "Данные пользователя, если поле пустое, то оно не меняется"
// @Success 200 {object} entities.UserProfile
// @Failure 500 {object} httperrors.HttpError
// @Router /api/profile/{id}/edit [post]
func (h *ProfileHandler) Edit(ctx context.Context, requestData entities.ProfileRequest) (entities.UserProfile, error) {
	username := requestData.Username
	bio := requestData.Bio
	avatar := requestData.Avatar

	pathParams := httputils.GetPathParamsFromCtx(ctx)
	idStr := pathParams["id"]
	userID, err := strconv.Atoi(idStr)
	if err != nil {
		return entities.UserProfile{}, err
	}

	newData := entities.UserProfile{
		UserID:   userID,
		Username: username,
		Bio:      bio,
		Avatar:   avatar,
	}
	err = h.userProfileUseCase.EditUserProfile(newData)
	if err != nil {
		return entities.UserProfile{}, err
	}
	return h.userProfileUseCase.GetUserProfile(userID)
}

// UploadFile TODO: нужно будет убрать это инженерное решение (костыль) после фикса обертки
func (h *ProfileHandler) UploadFile(w http.ResponseWriter, r *http.Request) {
	logger := logger.Logger()

	userID, err := strconv.Atoi(wrapper.GetPathParams(r)["id"])
	if err != nil {
		logger.Error("Handler error", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	path, err := delivery.SaveFile(r)
	if err != nil {
		logger.Error("Handler error", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	profile, err := h.userProfileUseCase.EditUserProfileAvatar(userID, path)
	if err != nil {
		logger.Error("Handler error", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	rawJSON, err := json.Marshal(profile)
	if err != nil {
		logger.Error("Error encoding response", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(rawJSON)
}
