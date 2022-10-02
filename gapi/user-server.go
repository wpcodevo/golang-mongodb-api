package gapi

import (
	"github.com/wpcodevo/golang-mongodb/config"
	"github.com/wpcodevo/golang-mongodb/pb"
	"github.com/wpcodevo/golang-mongodb/services"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserServer struct {
	pb.UnimplementedUserServiceServer
	config         config.Config
	userService    services.UserService
	userCollection *mongo.Collection
}

func NewGrpcUserServer(config config.Config, userService services.UserService, userCollection *mongo.Collection) (*UserServer, error) {
	userServer := &UserServer{
		config:         config,
		userService:    userService,
		userCollection: userCollection,
	}

	return userServer, nil
}
