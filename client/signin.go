package client

import (
	"context"
	"examplegrpcgin/pb"
	"fmt"
	"log"
	"time"

	"google.golang.org/grpc"
)

type SignInUserClient struct {
	service pb.AuthServiceClient
}

func NewSignInUserClient(conn *grpc.ClientConn) *SignInUserClient {
	service := pb.NewAuthServiceClient(conn)

	return &SignInUserClient{service}
}

func (signInUserClient *SignInUserClient) SignInUser(credentials *pb.SignInUserInput) {

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	res, err := signInUserClient.service.SignInUser(ctx, credentials)

	if err != nil {
		log.Fatalf("SignInUser: %v", err)
	}

	fmt.Println(res)
}
