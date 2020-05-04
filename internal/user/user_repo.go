package user

import "alexey-ershkov/alexey-ershkov-DB.git/internal/models"

type Repository interface {
	InsertInto(user *models.User) error
	GetByNickname(user *models.User) error
	GetByNicknameOrEmail(user *models.User) ([]models.User, error)
	Update(user *models.User) error
}
