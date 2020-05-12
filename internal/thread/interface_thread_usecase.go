package thread

import "alexey-ershkov/alexey-ershkov-DB.git/internal/models"

type Usecase interface {
	CreateThread(thread *models.Thread) error
	GetThreadInfo(thread *models.Thread) error
	CreateVote(thread *models.Thread, vote *models.Vote) error
	UpdateThread(thread *models.Thread) error
	GetThreadPosts(thread *models.Thread, desc, sort, limit, since string) ([]models.Post, error)
}
