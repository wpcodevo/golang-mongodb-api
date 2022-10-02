package gapi

import (
	"context"
	"strings"

	"github.com/wpcodevo/golang-mongodb/models"
	"github.com/wpcodevo/golang-mongodb/pb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

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
