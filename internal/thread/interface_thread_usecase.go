package thread

import "alexey-ershkov/alexey-ershkov-DB.git/internal/models"

type Usecase interface {
	CreateThread(thread *models.Thread) error
	GetThreadInfo(thread *models.Thread) error
}
