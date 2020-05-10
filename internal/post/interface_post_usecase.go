package post

import "alexey-ershkov/alexey-ershkov-DB.git/internal/models"

type Usecase interface {
	CreatePosts(posts []*models.Post, thread *models.Thread) error
	GetPost(post *models.Post) error
	UpdatePost(post *models.Post) error
}
