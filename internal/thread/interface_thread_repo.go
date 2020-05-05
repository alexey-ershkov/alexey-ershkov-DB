package thread

import "alexey-ershkov/alexey-ershkov-DB.git/internal/models"

type Repository interface {
	InsertInto(thread *models.Thread) error
	Get(thread *models.Thread) error
}
