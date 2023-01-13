package main

import (
	"context"
	"examplegrpcgin/config"
	"examplegrpcgin/gapi"
	"examplegrpcgin/pb"
	"examplegrpcgin/services"
	"log"
	"net"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"go.mongodb.org/mongo-driver/mongo"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var (
	server      *gin.Engine
	ctx         context.Context
	mongoclient *mongo.Client
	redisclient *redis.Client

	userService services.UserService

	authCollection *mongo.Collection
	authService    services.AuthService

	// 👇 Create the Post Variables
	postService    services.PostService
	postCollection *mongo.Collection
)

func main() {
	config, err := config.LoadConfig(".")

	if err != nil {
		log.Fatal("Could not load config", err)
	}
	defer mongoclient.Disconnect(ctx)
	startGrpcServer(config)

}

func startGrpcServer(config config.Config) {
	authServer, err := gapi.NewGrpcAuthServer(config, authService, userService, authCollection)
	if err != nil {
		log.Fatal("cannot create grpc authServer: ", err)
	}

	userServer, err := gapi.NewGrpcUserServer(config, userService, authCollection)
	if err != nil {
		log.Fatal("cannot create grpc userServer: ", err)
	}

	postServer, err := gapi.NewGrpcPostServer(postCollection, postService)
	if err != nil {
		log.Fatal("cannot create grpc postServer: ", err)
	}

	grpcServer := grpc.NewServer()

	pb.RegisterAuthServiceServer(grpcServer, authServer)
	pb.RegisterUserServiceServer(grpcServer, userServer)
	// 👇 Register the Post gRPC service
	pb.RegisterPostServiceServer(grpcServer, postServer)
	reflection.Register(grpcServer)

	listener, err := net.Listen("tcp", config.GrpcServerAddress)
	if err != nil {
		log.Fatal("cannot create grpc server: ", err)
	}

	log.Printf("start gRPC server on %s", listener.Addr().String())
	err = grpcServer.Serve(listener)
	if err != nil {
		log.Fatal("cannot create grpc server: ", err)
	}
}
