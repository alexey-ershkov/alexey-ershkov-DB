package usecase

import (
	"alexey-ershkov/alexey-ershkov-DB.git/internal/models"
	"alexey-ershkov/alexey-ershkov-DB.git/internal/thread"
	"alexey-ershkov/alexey-ershkov-DB.git/internal/tools"
	"github.com/sirupsen/logrus"
)

type Usecase struct {
	repo thread.Repository
}

func NewThreadUsecase(r thread.Repository) thread.Usecase {
	return &Usecase{
		repo: r,
	}
}

func (tUC *Usecase) CreateThread(th *models.Thread) error {
	if err := tUC.repo.InsertInto(th); err != nil {
		if err := tUC.repo.GetBySlug(th); err != nil {
			return tools.UserNotExist
		}
		return tools.ThreadExist
	}
	if err := tUC.repo.Get(th); err != nil {
		return err
	}
	return nil
}

func (tUC *Usecase) GetThreadInfo(th *models.Thread) error {
	if err := tUC.repo.Get(th); err != nil {
		logrus.Warn("thread doesn't exist")
		return tools.ThreadNotExist
	}
	return nil
}
