package forum

import "alexey-ershkov/alexey-ershkov-DB.git/internal/models"

type Usecase interface {
	CreateForum(forum *models.Forum) error
	GetForum(forum *models.Forum) error
	GetForumThreads(forum *models.Forum, desc, limit, since string) ([]models.Thread, error)
}
