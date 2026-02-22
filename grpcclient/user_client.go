package grpcclient

import (
	"context"
	"time"

	userpb "github.com/jrmygp/contracts/proto/userpb"
	"google.golang.org/grpc"
)

type UserClient struct {
	client userpb.UserServiceClient
}

func NewUserClient() (*UserClient, *grpc.ClientConn, error) {
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		return nil, nil, err
	}

	client := userpb.NewUserServiceClient(conn)

	return &UserClient{client: client}, conn, nil
}

func (u *UserClient) GetUserByID(id int32) (*userpb.GetUserResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return u.client.GetUserByID(ctx, &userpb.GetUserRequest{
		Id: id,
	})
}

func (u *UserClient) RefundBalance(
	userID int32,
	amount int32,
) (*userpb.EditUserResponse, error) {

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return u.client.EditUser(ctx, &userpb.EditUserRequest{
		UserId:          userID,
		IncomingBalance: amount,
	})
}
