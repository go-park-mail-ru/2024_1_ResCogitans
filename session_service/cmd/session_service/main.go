package main

import (
	"log"
	"net"

	"github.com/go-park-mail-ru/2024_1_ResCogitans/session_service/gen"
	"github.com/go-park-mail-ru/2024_1_ResCogitans/session_service/internal/config"
	"github.com/go-park-mail-ru/2024_1_ResCogitans/session_service/internal/database"
	"github.com/go-park-mail-ru/2024_1_ResCogitans/session_service/internal/session"
	"github.com/go-park-mail-ru/2024_1_ResCogitans/session_service/internal/storage"
	"google.golang.org/grpc"
)

func main() {
	cfg, err := config.LoadConfig()
	lis, err := net.Listen(cfg.Server.Protocol, cfg.Server.Port)
	if err != nil {
		log.Fatalln("can't listen port", err)
		return
	}

	server := grpc.NewServer()

	redisDB, err := database.GetSessionRedis()
	if err != nil {
		log.Fatalln(err)
		return
	}
	sessionStorage := storage.NewSessionStorage(redisDB)

	gen.RegisterSessionServiceServer(server, session.NewSessionManager(sessionStorage))

	err = server.Serve(lis)
	if err != nil {
		return
	}
}
