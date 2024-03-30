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
	if false {
		signUpUserClient := client.NewSignUpUserClient(conn)
		newUser := &pb.SignUpUserInput{
			Name:            "Jane Smith",
			Email:           "janesmith@gmail.com",
			Password:        "password123",
			PasswordConfirm: "password123",
		}
		signUpUserClient.SignUpUser(newUser)
	}

	// Sign In
	if false {
		signInUserClient := client.NewSignInUserClient(conn)

		credentials := &pb.SignInUserInput{
			Email:    "janesmith@gmail.com",
			Password: "password123",
		}
		signInUserClient.SignInUser(credentials)
	}

	// Get Me
	if false {

		getMeClient := client.NewGetMeClient(conn)
		id := &pb.GetMeRequest{
			Id: "628cffb91e50302d360c1a2c",
		}
		getMeClient.GetMeUser(id)

	}

	// List Posts
	if false {
		listPostsClient := client.NewListPostsClient(conn)

		var page int64 = 1
		var limit int64 = 10
		args := &pb.GetPostsRequest{
			Page:  &page,
			Limit: &limit,
		}

		listPostsClient.ListPosts(args)
	}

	// Create Post
	if true {
		createPostClient := client.NewCreatePostClient(conn)

		args := &pb.CreatePostRequest{
			Title:   "My second gRPC post with joy",
			Content: "It's always good to learn new technologies",
			User:    "62908e0a42a608d5aeae2f64",
			Image:   "default.png",
		}

		createPostClient.CreatePost(args)
	}

	// Update Post
	if false {
		updatePostClient := client.NewUpdatePostClient(conn)

		title := "My new updated title for my blog"
		args := &pb.UpdatePostRequest{
			Id:    "629169e00a6c7cfd24e2129d",
			Title: &title,
		}

		updatePostClient.UpdatePost(args)
	}

	// Get Post
	if false {
		getPostClient := client.NewGetPostClient(conn)

		args := &pb.PostRequest{
			Id: "629169e00a6c7cfd24e2129d",
		}

		getPostClient.GetPost(args)
	}

	// Delete Post
	if false {
		deletePostClient := client.NewDeletePostClient(conn)

		args := &pb.PostRequest{
			Id: "629147ff3c92aed11d49394b",
		}

		deletePostClient.DeletePost(args)
	}
}
