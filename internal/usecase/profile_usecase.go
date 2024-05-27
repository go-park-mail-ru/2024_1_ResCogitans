package usecase

import (
	"context"

	"github.com/go-park-mail-ru/2024_1_ResCogitans/internal/entities"
	"github.com/go-park-mail-ru/2024_1_ResCogitans/internal/storage/postgres/user"
)

type ProfileUseCaseInterface interface {
	GetUserProfile(ctx context.Context, userID int) (entities.UserProfile, error)
	EditUserProfile(ctx context.Context, newData entities.UserProfile) error
	EditUserProfileAvatar(ctx context.Context, userID int, avatar string) (entities.UserProfile, error)
}

type ProfileUseCase struct {
	ProfileStorage *user.UserProfileStorage
}

func NewProfileUseCase(storage *user.UserProfileStorage) *ProfileUseCase {
	return &ProfileUseCase{
		ProfileStorage: storage,
	}
}

func (pu *ProfileUseCase) GetUserProfile(ctx context.Context, userID int) (entities.UserProfile, error) {
	return pu.ProfileStorage.GetUserProfileByID(ctx, userID)
}

func (pu *ProfileUseCase) EditUserProfile(ctx context.Context, newData entities.UserProfile) error {
	if newData.Username != "" {
		err := pu.ProfileStorage.EditUsername(ctx, newData.UserID, newData.Username)
		if err != nil {
			return err
		}
	}

	if newData.Bio != "" {
		err := pu.ProfileStorage.EditUserBio(ctx, newData.UserID, newData.Bio)
		if err != nil {
			return err
		}
	}

	if newData.Avatar != "" {
		err := pu.ProfileStorage.EditUserAvatar(ctx, newData.UserID, newData.Avatar)
		if err != nil {
			return err
		}
	}
	return nil
}

func (pu *ProfileUseCase) EditUserProfileAvatar(ctx context.Context, userID int, avatar string) (entities.UserProfile, error) {
	err := pu.ProfileStorage.EditUserAvatar(ctx, userID, avatar)
	if err != nil {
		return entities.UserProfile{}, err
	}
	return pu.ProfileStorage.GetUserProfileByID(ctx, userID)
}
