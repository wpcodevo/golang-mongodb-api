package client

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/wpcodevo/golang-mongodb/pb"
	"google.golang.org/grpc"
)

type GetMeClient struct {
	service pb.UserServiceClient
}

func NewGetMeClient(conn *grpc.ClientConn) *GetMeClient {
	service := pb.NewUserServiceClient(conn)

	return &GetMeClient{service}
}

func (getMeClient *GetMeClient) GetMeUser(credentials *pb.GetMeRequest) {

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(time.Millisecond*5000))
	defer cancel()

	res, err := getMeClient.service.GetMe(ctx, credentials)

	if err != nil {
		log.Fatalf("GeMe: %v", err)
	}

	fmt.Println(res)
}
