package main

import (
	"log"

	"github.com/wpcodevo/golang-mongodb/client"
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

	if false {
		signUpUserClient := client.NewSignUpUserClient(conn)
		newUser := &pb.SignUpUserInput{
			Name:            "Micheal Smith",
			Email:           "michealmith@gmail.com",
			Password:        "password123",
			PasswordConfirm: "password123",
		}
		signUpUserClient.SignUpUser(newUser)
	}

	getMeClient := client.NewGetMeClient(conn)
	id := &pb.GetMeRequest{
		Id: "628cffb91e50302d360c1a2c",
	}
	getMeClient.GetMeUser(id)

}
