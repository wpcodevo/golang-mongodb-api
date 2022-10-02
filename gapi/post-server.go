package gapi

import (
	"github.com/wpcodevo/golang-mongodb/pb"
	"github.com/wpcodevo/golang-mongodb/services"
	"go.mongodb.org/mongo-driver/mongo"
)

type PostServer struct {
	pb.UnimplementedPostServiceServer
	postCollection *mongo.Collection
	postService    services.PostService
}

func NewGrpcPostServer(postCollection *mongo.Collection, postService services.PostService) (*PostServer, error) {
	postServer := &PostServer{
		postCollection: postCollection,
		postService:    postService,
	}

	return postServer, nil
}
