package main

import (
	"fmt"
	"log"
	"net"

	"github.com/go-park-mail-ru/2024_1_ResCogitans/internal/service/gen"
	"github.com/go-park-mail-ru/2024_1_ResCogitans/session_service/database"
	"github.com/go-park-mail-ru/2024_1_ResCogitans/session_service/session"
	"google.golang.org/grpc"
)

func main() {
	lis, err := net.Listen("tcp", ":8081")
	if err != nil {
		log.Fatalln("can't listen port", err)
	}

	server := grpc.NewServer()

	pdb, err := database.GetSessionRedis()
	if err != nil {
		log.Fatalln(err)
	}

	gen.RegisterSessionServiceServer(server, session.NewSessionManager(pdb))

	fmt.Println("starting server at :8081")
	err = server.Serve(lis)
	if err != nil {
		return
	}
}
