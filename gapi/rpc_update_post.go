package gapi

import (
	"context"
	"strings"
	"time"

	"github.com/wpcodevo/golang-mongodb/models"
	"github.com/wpcodevo/golang-mongodb/pb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

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
