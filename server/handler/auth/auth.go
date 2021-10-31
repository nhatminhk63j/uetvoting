package auth

import (
	"context"

	pb "github.com/nhatminhk63j/uetvoting/pb/auth/v1"
	"github.com/nhatminhk63j/uetvoting/pkg/auth"
)

type ServiceServer struct {
	*pb.UnimplementedAuthServiceServer

	authSvc auth.Service
}

// NewServiceServer ...
func NewServiceServer(authSvc auth.Service) *ServiceServer {
	return &ServiceServer{
		authSvc: authSvc,
	}
}

// Login ...
func (s *ServiceServer) Login(ctx context.Context, in *pb.LoginRequest) (*pb.LoginResponse, error) {
	accessToken, err := s.authSvc.Login(ctx, in.IdToken)
	if err != nil {
		return nil, err
	}
	return &pb.LoginResponse{
		AccessToken: accessToken,
	}, nil
}
