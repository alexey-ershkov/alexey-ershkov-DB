package forum

import "alexey-ershkov/alexey-ershkov-DB.git/internal/models"

type Repository interface {
	InsertInto(forum *models.Forum) error
	GetBySlug(forum *models.Forum) error
}
