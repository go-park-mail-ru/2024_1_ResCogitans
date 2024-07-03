package main

import (
	"fmt"
	"log"

	"github.com/go-park-mail-ru/2024_1_ResCogitans/internal/config"
	"github.com/go-park-mail-ru/2024_1_ResCogitans/internal/delivery/initialization"
	"github.com/go-park-mail-ru/2024_1_ResCogitans/internal/delivery/server"
	"github.com/go-park-mail-ru/2024_1_ResCogitans/internal/usecase"
	"github.com/go-park-mail-ru/2024_1_ResCogitans/router"
	"github.com/go-park-mail-ru/2024_1_ResCogitans/utils/logger"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	logger := logger.Logger()
	cfg, err := config.LoadConfig()
	if err != nil {
		logger.Error("Failed to load config", "error", err)
		return
	}
	logger.Info("Start config")

	conn, err := grpc.Dial(fmt.Sprintf("%s:%d", cfg.SessionService.Host, cfg.SessionService.Port), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logger.Error("did not connect: %v", err)
	}
	defer func(conn *grpc.ClientConn) {
		err := conn.Close()
		if err != nil {
			log.Fatalf("could not close connection: %v", err)
		}
	}(conn)

	postgresDB, CSRFDB, err := initialization.DataBaseInitialization()
	if err != nil {
		logger.Error("DataBase initialization error", "error", err)
		return
	}

	SessionUseCase := usecase.NewSessionUseCase(conn)

	storages := initialization.StorageInit(postgresDB, CSRFDB)
	usecases := initialization.UseCaseInit(storages, SessionUseCase)
	handlers := initialization.HandlerInit(usecases)

	router := router.SetupRouter(cfg, handlers)

	if err := server.StartServer(router, cfg); err != nil {
		logger.Error("Failed to start server", "error", err)
		return
	}
}
