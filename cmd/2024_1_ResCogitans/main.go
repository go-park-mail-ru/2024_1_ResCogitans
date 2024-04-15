package main

import (
	"context"
	"fmt"

	"github.com/go-park-mail-ru/2024_1_ResCogitans/internal/config"
	"github.com/go-park-mail-ru/2024_1_ResCogitans/internal/http-server/handlers/authorization"
	"github.com/go-park-mail-ru/2024_1_ResCogitans/internal/http-server/handlers/deactivation"
	"github.com/go-park-mail-ru/2024_1_ResCogitans/internal/http-server/handlers/registration"
	"github.com/go-park-mail-ru/2024_1_ResCogitans/internal/http-server/handlers/updateUserData"
	"github.com/go-park-mail-ru/2024_1_ResCogitans/internal/http-server/server"
	"github.com/go-park-mail-ru/2024_1_ResCogitans/internal/storage/session"
	"github.com/go-park-mail-ru/2024_1_ResCogitans/internal/storage/user"
	"github.com/go-park-mail-ru/2024_1_ResCogitans/internal/usecase"
	"github.com/go-park-mail-ru/2024_1_ResCogitans/router"
	"github.com/go-park-mail-ru/2024_1_ResCogitans/utils/app"
	"github.com/go-park-mail-ru/2024_1_ResCogitans/utils/logger"
	"github.com/go-park-mail-ru/2024_1_ResCogitans/utils/middle"
	"github.com/go-redis/redis/v8"
)

func main() {
	logger := logger.Logger()
	cfg, err := config.LoadConfig()
	if err != nil {
		logger.Error("Failed to load config", "error", err)
		return
	}

	rdb := redis.NewClient(&redis.Options{
		Addr:     "redis-13041.c302.asia-northeast1-1.gce.cloud.redislabs.com:13041",
		Username: "default",
		Password: "Hwsuxke8YC8vT6E2jOKd7lTK6cPEvq5I", // Замените на ваш пароль
		DB:       0,                                  // Обычно используется 0 для первой базы данных
	})

	// Используем контекст с отменой по умолчанию для выполнения операций Redis
	ctx := context.Background()

	// Проверяем подключение
	pong, err := rdb.Ping(ctx).Result()
	if err != nil {
		panic(err)
	}
	fmt.Println(pong, "response from Redis")

	redisStorage := session.NewRedisStorage("redis-13041.c302.asia-northeast1-1.gce.cloud.redislabs.com:13041", "default", "Hwsuxke8YC8vT6E2jOKd7lTK6cPEvq5I", 0)
	sessionUseCase := usecase.NewSessionUseCase(redisStorage)

	userStorage := user.NewUserStorage()
	userUseCase := usecase.NewUserUseCase(userStorage)

	authHandler := authorization.NewAuthorizationHandler(sessionUseCase, userUseCase)
	regHandler := registration.NewRegistrationHandler(sessionUseCase, userUseCase)
	updateHandler := updateUserData.NewUpdateDataHandler(sessionUseCase, userUseCase)
	deactivationHandler := deactivation.NewDeactivationHandler(sessionUseCase, userUseCase)
	authMiddleware := middle.NewAuthMiddleware(sessionUseCase)

	app := app.NewApp(authHandler, regHandler, updateHandler, deactivationHandler, authMiddleware)

	logger.Info("Start config", "config", cfg)

	router := router.SetupRouter(app)

	if err := server.StartServer(router, cfg); err != nil {
		logger.Error("Failed to start server", "error", err)
		return
	}
}
