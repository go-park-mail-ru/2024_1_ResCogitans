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

func (h *ProfileHandler) Get(ctx context.Context, _ entities.UserProfile) (entities.UserProfile, error) {
	pathParams := httputils.GetPathParamsFromCtx(ctx)
	idStr := pathParams["id"]
	userID, err := strconv.Atoi(idStr)
	if err != nil {
		return entities.UserProfile{}, err
	}
	return h.userProfileUseCase.GetUserProfile(userID)
}

func (h *ProfileHandler) Edit(ctx context.Context, requestData entities.UserProfile) (entities.UserProfile, error) {
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
