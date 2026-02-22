package chirps

import (
	"chirpstream/internal/pb"
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type UserClient struct {
	client pb.UserServiceClient
}

func NewUserClient(target string) (*UserClient, error) {
	conn, err := grpc.NewClient(target, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	return &UserClient{
		client: pb.NewUserServiceClient(conn),
	}, nil
}

func (c *UserClient) GetUser(ctx context.Context, id int) (*pb.GetUserByIdResponse, error) {
	req := &pb.GetUserByIdRequest{Id: int32(id)}

	return c.client.GetUserById(ctx, req)
}
