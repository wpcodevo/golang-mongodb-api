package services

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/wpcodevo/golang-mongodb/models"
	"github.com/wpcodevo/golang-mongodb/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type UserServiceImpl struct {
	collection *mongo.Collection
	ctx        context.Context
}

func NewUserServiceImpl(collection *mongo.Collection, ctx context.Context) UserService {
	return &UserServiceImpl{collection, ctx}
}

func (us *UserServiceImpl) FindUserById(id string) (*models.DBResponse, error) {
	oid, _ := primitive.ObjectIDFromHex(id)

	var user *models.DBResponse

	query := bson.M{"_id": oid}
	err := us.collection.FindOne(us.ctx, query).Decode(&user)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return &models.DBResponse{}, err
		}
		return nil, err
	}

	return user, nil
}

func (us *UserServiceImpl) FindUserByEmail(email string) (*models.DBResponse, error) {
	var user *models.DBResponse

	query := bson.M{"email": strings.ToLower(email)}
	err := us.collection.FindOne(us.ctx, query).Decode(&user)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return &models.DBResponse{}, err
		}
		return nil, err
	}

	return user, nil
}

func (uc *UserServiceImpl) UpdateUserById(id string, data *models.UpdateInput) (*models.DBResponse, error) {
	doc, err := utils.ToDoc(data)
	if err != nil {
		return &models.DBResponse{}, err
	}

	fmt.Println(data)

	obId, _ := primitive.ObjectIDFromHex(id)

	query := bson.D{{Key: "_id", Value: obId}}
	update := bson.D{{Key: "$set", Value: doc}}
	result := uc.collection.FindOneAndUpdate(uc.ctx, query, update, options.FindOneAndUpdate().SetReturnDocument(1))

	var updatedUser *models.DBResponse
	if err := result.Decode(&updatedUser); err != nil {
		return nil, errors.New("no document with that id exists")
	}

	return updatedUser, nil
}
