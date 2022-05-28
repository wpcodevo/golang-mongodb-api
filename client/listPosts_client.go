package client

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"

	"github.com/wpcodevo/golang-mongodb/pb"
	"google.golang.org/grpc"
)

type ListPostsClient struct {
	service pb.PostServiceClient
}

func NewListPostsClient(conn *grpc.ClientConn) *ListPostsClient {
	service := pb.NewPostServiceClient(conn)

	return &ListPostsClient{service}
}

func (listPostsClient *ListPostsClient) ListPosts(args *pb.GetPostsRequest) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(time.Millisecond*5000))
	defer cancel()

	stream, err := listPostsClient.service.GetPosts(ctx, args)
	if err != nil {
		log.Fatalf("ListPosts: %v", err)
	}

	for {
		res, err := stream.Recv()

		if err == io.EOF {
			break
		}

		if err != nil {
			log.Fatalf("ListPosts: %v", err)
		}

		fmt.Println(res)
	}

}
