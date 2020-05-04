package user

import "alexey-ershkov/alexey-ershkov-DB.git/internal/models"

type Repository interface {
	InsertInto(user *models.User) error
	GetUserByNickname(user *models.User) error
	GetUsersByNicknameOrEmail(user *models.User) ([]models.User, error)
}
