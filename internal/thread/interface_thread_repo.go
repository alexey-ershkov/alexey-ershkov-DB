package thread

import "alexey-ershkov/alexey-ershkov-DB.git/internal/models"

type Repository interface {
	InsertInto(thread *models.Thread) error
	GetCreated(thread *models.Thread) error
	GetBySlug(thread *models.Thread) error
	GetBySlugOrId(thread *models.Thread) error
	InsertIntoVotes(vote *models.Vote) error
	GetVotes(thread *models.Thread, vote *models.Vote) error
}
