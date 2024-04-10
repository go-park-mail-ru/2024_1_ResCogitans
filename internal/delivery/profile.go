package delivery

import (
	"context"
	"strconv"

	"github.com/go-park-mail-ru/2024_1_ResCogitans/internal/entities"
	"github.com/go-park-mail-ru/2024_1_ResCogitans/internal/http-server/server/db"
	userRep "github.com/go-park-mail-ru/2024_1_ResCogitans/internal/repository/postgres"
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
