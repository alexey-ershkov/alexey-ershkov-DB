package usecase

import (
	"alexey-ershkov/alexey-ershkov-DB.git/internal/forum"
	"alexey-ershkov/alexey-ershkov-DB.git/internal/models"
)

type Usecase struct {
	repo forum.Repository
}

func NewForumUsecase(r forum.Repository) forum.Usecase {
	return &Usecase{
		repo: r,
	}
}

func (uc *Usecase) CreateForum(f *models.Forum) error {
	err := uc.repo.InsertInto(f)
	if err != nil {
		return err
	}
	return nil
}

func (uc *Usecase) GetForum(f *models.Forum) error {
	err := uc.repo.GetBySlug(f)
	if err != nil {
		return err
	}
	return nil
}
