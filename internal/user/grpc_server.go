package user

import (
	"chirpstream/internal/pb"
	"context"
	"fmt"
)

type GRPCServer struct {
	pb.UnimplementedUserServiceServer
	service Service
}

func NewGRPCServer(service Service) *GRPCServer {
	return &GRPCServer{service: service}
}

func (s *GRPCServer) GetUserById(ctx context.Context, req *pb.GetUserByIdRequest) (*pb.GetUserByIdResponse, error) {
	user, err := s.service.GetUserById(ctx, int(req.Id))
	if err != nil {
		return nil, fmt.Errorf("User not found %w", err)
	}

	return &pb.GetUserByIdResponse{
		Id:       int32(user.ID),
		Name:     user.Name,
		Username: user.Username,
	}, nil
}
