package cmd

import (
	"log"
	"net"

	"github.com/go-park-mail-ru/2024_1_ResCogitans/internal/service/gen"
	"google.golang.org/grpc"
)

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	sessionService := gen.NewSessionService()

	gen.RegisterSessionServiceServer(grpcServer, sessionService)

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
