package services

import "github.com/wpcodevo/golang-mongodb/models"

type UserService interface {
	FindUserById(id string) (*models.DBResponse, error)
	FindUserByEmail(email string) (*models.DBResponse, error)
	UpdateUserById(id string, data *models.UpdateInput) (*models.DBResponse, error)
	UpdateOne(field string, value interface{}) (*models.DBResponse, error)
}
