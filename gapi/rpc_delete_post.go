package gapi

import (
	"context"
	"strings"

	"github.com/wpcodevo/golang-mongodb/pb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

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
