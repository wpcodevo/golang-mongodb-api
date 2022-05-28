package services

import "github.com/wpcodevo/golang-mongodb/models"

type AuthService interface {
	SignUpUser(*models.SignUpInput) (*models.UserDBResponse, error)
	SignInUser(*models.SignInInput) (*models.UserDBResponse, error)
}
