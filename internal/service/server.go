package service

import (
	"context"
)

type grpcServer struct {
	sessionService *SessionService
}

func (s *grpcServer) SaveSession(ctx context.Context, req *SaveSessionRequest) (*proto.SaveSessionResponse, error) {
}
