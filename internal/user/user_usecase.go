package user

import "alexey-ershkov/alexey-ershkov-DB.git/internal/models"

type Usecase interface {
	CreateUser(user *models.User) ([]models.User, error)
	GetUser(user *models.User) error
}
