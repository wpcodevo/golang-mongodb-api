package client

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/wpcodevo/golang-mongodb/pb"
	"google.golang.org/grpc"
)

type UpdatePostClient struct {
	service pb.PostServiceClient
}

func NewUpdatePostClient(conn *grpc.ClientConn) *UpdatePostClient {
	service := pb.NewPostServiceClient(conn)

	return &UpdatePostClient{service}
}

func (updatePostClient *UpdatePostClient) UpdatePost(args *pb.UpdatePostRequest) {

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(time.Millisecond*5000))
	defer cancel()

	res, err := updatePostClient.service.UpdatePost(ctx, args)

	if err != nil {
		log.Fatalf("UpdatePost: %v", err)
	}

	fmt.Println(res)
}
