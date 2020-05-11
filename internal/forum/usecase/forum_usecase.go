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

func (uc *Usecase) GetForumThreads(f *models.Forum, desc, limit, since string) ([]models.Thread, error) {
	ths, err := uc.repo.GetThreads(f, desc, limit, since)
	if err != nil {
		return nil, err
	}
	return ths, nil
}

func (uc *Usecase) GetForumUsers(f *models.Forum, desc, limit, since string) ([]models.User, error) {
	usr, err := uc.repo.GetUsers(f, desc, limit, since)
	if err != nil {
		return nil, err
	}
	return usr, nil
}
