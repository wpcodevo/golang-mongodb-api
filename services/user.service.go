package services

import "github.com/wpcodevo/golang-mongodb/models"

type UserService interface {
	FindUserById(id string) (*models.UserDBResponse, error)
	FindUserByEmail(email string) (*models.UserDBResponse, error)
	UpdateUserById(id string, data *models.UserUpdateInput) (*models.UserDBResponse, error)
}
