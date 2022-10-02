package gapi

import (
	"github.com/wpcodevo/golang-mongodb/pb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

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
