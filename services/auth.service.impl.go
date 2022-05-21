package services

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/wpcodevo/golang-mongodb/models"
	"github.com/wpcodevo/golang-mongodb/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type AuthServiceImpl struct {
	collection *mongo.Collection
	ctx        context.Context
}

func NewAuthService(collection *mongo.Collection, ctx context.Context) AuthService {
	return &AuthServiceImpl{collection, ctx}
}

func (uc *AuthServiceImpl) SignUpUser(user *models.SignUpInput) (*models.DBResponse, error) {
	user.CreatedAt = time.Now()
	user.UpdatedAt = user.CreatedAt
	user.Email = strings.ToLower(user.Email)
	user.PasswordConfirm = ""
	user.Verified = false
	user.Role = "user"

	hashedPassword, _ := utils.HashPassword(user.Password)
	user.Password = hashedPassword
	res, err := uc.collection.InsertOne(uc.ctx, &user)

	if err != nil {
		if er, ok := err.(mongo.WriteException); ok && er.WriteErrors[0].Code == 11000 {
			return nil, errors.New("user with that email already exist")
		}
		return nil, err
	}

	// Create a unique index for the email field
	opt := options.Index()
	opt.SetUnique(true)
	index := mongo.IndexModel{Keys: bson.M{"email": 1}, Options: opt}

	if _, err := uc.collection.Indexes().CreateOne(uc.ctx, index); err != nil {
		return nil, errors.New("could not create index for email")
	}

	var newUser *models.DBResponse
	query := bson.M{"_id": res.InsertedID}

	err = uc.collection.FindOne(uc.ctx, query).Decode(&newUser)
	if err != nil {
		return nil, err
	}

	return newUser, nil
}

func (uc *AuthServiceImpl) SignInUser(*models.SignInInput) (*models.DBResponse, error) {
	return nil, nil
}
