package profile

import (
	"context"
	"strconv"

	"github.com/go-park-mail-ru/2024_1_ResCogitans/internal/entities"
	"github.com/go-park-mail-ru/2024_1_ResCogitans/internal/usecase"
	"github.com/go-park-mail-ru/2024_1_ResCogitans/utils/httputils"
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
