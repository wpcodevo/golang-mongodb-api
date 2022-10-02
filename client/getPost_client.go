package client

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/wpcodevo/golang-mongodb/pb"
	"google.golang.org/grpc"
)

type GetPostClient struct {
	service pb.PostServiceClient
}

func NewGetPostClient(conn *grpc.ClientConn) *GetPostClient {
	service := pb.NewPostServiceClient(conn)

	return &GetPostClient{service}
}

func (getPostClient *GetPostClient) GetPost(args *pb.PostRequest) {

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(time.Millisecond*5000))
	defer cancel()

	res, err := getPostClient.service.GetPost(ctx, args)

	if err != nil {
		log.Fatalf("GetPost: %v", err)
	}

	fmt.Println(res)
}
