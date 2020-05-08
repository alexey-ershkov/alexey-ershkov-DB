package post

import "alexey-ershkov/alexey-ershkov-DB.git/internal/models"

type Repository interface {
	InsertInto(post []*models.Post) error
	GetById(post *models.Post) error
}
