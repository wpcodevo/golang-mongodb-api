package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/wpcodevo/golang-mongodb/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	address = "0.0.0.0:8080"
)

func main() {
	conn, err := grpc.Dial(address, grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock())

	if err != nil {
		log.Fatalf("failed to connect: %v", err)
	}

	defer conn.Close()

	client := pb.NewAuthServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(time.Millisecond*5000))
	defer cancel()

	newUser := &pb.SignUpUserInput{
		Name:            "James Smith",
		Email:           "jamesmith@gmail.com",
		Password:        "password123",
		PasswordConfirm: "password123",
	}

	res, err := client.SignUpUser(ctx, newUser)
	if err != nil {
		log.Fatalf("SignUpUser: %v", err)
	}

	fmt.Println(res)

}
