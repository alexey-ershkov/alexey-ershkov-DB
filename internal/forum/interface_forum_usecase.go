package forum

import "alexey-ershkov/alexey-ershkov-DB.git/internal/models"

type Usecase interface {
	CreateForum(forum *models.Forum) error
	GetForum(forum *models.Forum) error
}
