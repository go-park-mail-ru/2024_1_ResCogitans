package usecase

import (
	"github.com/go-park-mail-ru/2024_1_ResCogitans/internal/entities"
	storage "github.com/go-park-mail-ru/2024_1_ResCogitans/internal/storage/storage_interfaces"
)

type ProfileUseCaseInterface interface {
	GetUserProfile(userID int) (entities.UserProfile, error)
	EditUserProfile(newData entities.UserProfile) error
}

type ProfileUseCase struct {
	ProfileStorage storage.UserProfileStorageInterface
}

func NewProfileUseCase(storage storage.UserProfileStorageInterface) ProfileUseCaseInterface {
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
