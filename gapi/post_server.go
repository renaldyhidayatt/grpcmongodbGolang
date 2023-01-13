package gapi

import (
	"context"
	"examplegrpcgin/models"
	"examplegrpcgin/pb"
	"examplegrpcgin/services"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type PostServer struct {
	pb.UnimplementedPostServiceServer
	postCollection *mongo.Collection
	postService    services.PostService
}

func NewGrpcPostServer(postCollection *mongo.Collection, postService services.PostService) (*PostServer, error) {
	postServer := &PostServer{
		postCollection: postCollection,
		postService:    postService,
	}

	return postServer, nil
}

func (postServer *PostServer) GetPost(ctx context.Context, req *pb.PostRequest) (*pb.PostResponse, error) {
	postId := req.GetId()

	post, err := postServer.postService.FindPostById(postId)
	if err != nil {
		if strings.Contains(err.Error(), "Id exists") {
			return nil, status.Errorf(codes.NotFound, err.Error())

		}
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	res := &pb.PostResponse{
		Post: &pb.Post{
			Id:        post.Id.Hex(),
			Title:     post.Title,
			Content:   post.Content,
			Image:     post.Image,
			User:      post.User,
			CreatedAt: timestamppb.New(post.CreateAt),
			UpdatedAt: timestamppb.New(post.UpdatedAt),
		},
	}
	return res, nil
}

func (postServer *PostServer) GetPosts(req *pb.GetPostsRequest, stream pb.PostService_GetPostsServer) error {
	var page = req.GetPage()
	var limit = req.GetLimit()

	posts, err := postServer.postService.FindPosts(int(page), int(limit))
	if err != nil {
		return status.Errorf(codes.Internal, err.Error())
	}

	for _, post := range posts {
		stream.Send(&pb.Post{
			Id:        post.Id.Hex(),
			Title:     post.Title,
			Content:   post.Content,
			Image:     post.Image,
			CreatedAt: timestamppb.New(post.CreateAt),
			UpdatedAt: timestamppb.New(post.UpdatedAt),
		})
	}

	return nil
}

func (postServer *PostServer) CreatePost(ctx context.Context, req *pb.CreatePostRequest) (*pb.PostResponse, error) {

	post := &models.CreatePostRequest{
		Title:   req.GetTitle(),
		Content: req.GetContent(),
		Image:   req.GetImage(),
		User:    req.GetUser(),
	}

	newPost, err := postServer.postService.CreatePost(post)

	if err != nil {
		if strings.Contains(err.Error(), "title already exists") {
			return nil, status.Errorf(codes.AlreadyExists, err.Error())
		}

		return nil, status.Errorf(codes.Internal, err.Error())
	}

	res := &pb.PostResponse{
		Post: &pb.Post{
			Id:        newPost.Id.Hex(),
			Title:     newPost.Title,
			Content:   newPost.Content,
			User:      newPost.User,
			CreatedAt: timestamppb.New(newPost.CreateAt),
			UpdatedAt: timestamppb.New(newPost.UpdatedAt),
		},
	}
	return res, nil
}

func (postServer *PostServer) UpdatePost(ctx context.Context, req *pb.UpdatePostRequest) (*pb.PostResponse, error) {
	postId := req.GetId()

	post := &models.UpdatePost{
		Title:     req.GetTitle(),
		Content:   req.GetContent(),
		Image:     req.GetImage(),
		User:      req.GetUser(),
		UpdatedAt: time.Now(),
	}

	updatedPost, err := postServer.postService.UpdatePost(postId, post)

	if err != nil {
		if strings.Contains(err.Error(), "Id exists") {
			return nil, status.Errorf(codes.NotFound, err.Error())
		}
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	res := &pb.PostResponse{
		Post: &pb.Post{
			Id:        updatedPost.Id.Hex(),
			Title:     updatedPost.Title,
			Content:   updatedPost.Content,
			Image:     updatedPost.Image,
			User:      updatedPost.User,
			CreatedAt: timestamppb.New(updatedPost.CreateAt),
			UpdatedAt: timestamppb.New(updatedPost.UpdatedAt),
		},
	}
	return res, nil
}

func (postServer *PostServer) DeletePost(ctx context.Context, req *pb.PostRequest) (*pb.DeletePostResponse, error) {
	postId := req.GetId()

	if err := postServer.postService.DeletePost(postId); err != nil {
		if strings.Contains(err.Error(), "Id exists") {
			return nil, status.Errorf(codes.NotFound, err.Error())
		}
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	res := &pb.DeletePostResponse{
		Success: true,
	}

	return res, nil
}
