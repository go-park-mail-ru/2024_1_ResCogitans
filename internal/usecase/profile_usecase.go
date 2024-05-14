package usecase

import (
	"github.com/go-park-mail-ru/2024_1_ResCogitans/internal/entities"
	"github.com/go-park-mail-ru/2024_1_ResCogitans/internal/storage/postgres/user"
)

type ProfileUseCaseInterface interface {
	GetUserProfile(userID int) (entities.UserProfile, error)
	EditUserProfile(newData entities.UserProfile) error
	EditUserProfileAvatar(userID int, avatar string) (entities.UserProfile, error)
}

type ProfileUseCase struct {
	ProfileStorage *user.UserProfileStorage
}

func NewProfileUseCase(storage *user.UserProfileStorage) *ProfileUseCase {
	return &ProfileUseCase{
		ProfileStorage: storage,
	}
}

func (pu *ProfileUseCase) GetUserProfile(userID int) (entities.UserProfile, error) {
	return pu.ProfileStorage.GetUserProfileByID(userID)
}

func (pu *ProfileUseCase) EditUserProfile(newData entities.UserProfile) error {
	if newData.Username != "" {
		err := pu.ProfileStorage.EditUsername(newData.UserID, newData.Username)
		if err != nil {
			return err
		}
	}

	if newData.Bio != "" {
		err := pu.ProfileStorage.EditUserBio(newData.UserID, newData.Bio)
		if err != nil {
			return err
		}
	}

	if newData.Avatar != "" {
		err := pu.ProfileStorage.EditUserAvatar(newData.UserID, newData.Avatar)
		if err != nil {
			return err
		}
	}
	return nil
}

func (pu *ProfileUseCase) EditUserProfileAvatar(userID int, avatar string) (entities.UserProfile, error) {
	err := pu.ProfileStorage.EditUserAvatar(userID, avatar)
	if err != nil {
		return entities.UserProfile{}, err
	}
	return pu.ProfileStorage.GetUserProfileByID(userID)
}
