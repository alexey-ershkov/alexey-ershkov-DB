package forum

import (
	"alexey-ershkov/alexey-ershkov-DB.git/internal/models"
)

type Repository interface {
	InsertInto(forum *models.Forum) error
	GetBySlug(forum *models.Forum) error
	GetThreads(forum *models.Forum, desc, limit, since string) ([]models.Thread, error)
	GetUsers(forum *models.Forum, desc, limit, since string) ([]models.User, error)
}