package client

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/wpcodevo/golang-mongodb/pb"
	"google.golang.org/grpc"
)

type DeletePostClient struct {
	service pb.PostServiceClient
}

func NewDeletePostClient(conn *grpc.ClientConn) *DeletePostClient {
	service := pb.NewPostServiceClient(conn)

	return &DeletePostClient{service}
}

func (deletePostClient *DeletePostClient) DeletePost(args *pb.PostRequest) {

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(time.Millisecond*5000))
	defer cancel()

	_, err := deletePostClient.service.DeletePost(ctx, args)

	if err != nil {
		log.Fatalf("DeletePost: %v", err)
	}

	fmt.Println("Post deleted successfully")
}
