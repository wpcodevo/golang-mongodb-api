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

	// Sign Up
	if 1 > 2 {
		signUpUserClient := client.NewSignUpUserClient(conn)
		newUser := &pb.SignUpUserInput{
			Name:            "Micheal Smith",
			Email:           "michealmith@gmail.com",
			Password:        "password123",
			PasswordConfirm: "password123",
		}
		signUpUserClient.SignUpUser(newUser)
	}

	// Sign In
	if 2 > 1 {
		signInUserClient := client.NewSignInUserClient(conn)

		credentials := &pb.SignInUserInput{
			Email:    "michealmith@gmail.com",
			Password: "password123",
		}
		signInUserClient.SignInUser(credentials)
	}

	// Get Me
	if 3 < 1 {

		getMeClient := client.NewGetMeClient(conn)
		id := &pb.GetMeRequest{
			Id: "628cffb91e50302d360c1a2c",
		}
		getMeClient.GetMeUser(id)

	}

}
