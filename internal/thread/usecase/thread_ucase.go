package usecase

import (
	"alexey-ershkov/alexey-ershkov-DB.git/internal/models"
	"alexey-ershkov/alexey-ershkov-DB.git/internal/thread"
	"alexey-ershkov/alexey-ershkov-DB.git/internal/tools"
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
	if err := tUC.repo.GetCreated(th); err != nil {
		return err
	}
	return nil
}

func (tUC *Usecase) GetThreadInfo(th *models.Thread) error {
	if err := tUC.repo.GetBySlugOrId(th); err != nil {
		//logrus.Warn("thread doesn't exist")
		return tools.ThreadNotExist
	}
	return nil
}

func (tUC *Usecase) CreateVote(th *models.Thread, v *models.Vote) error {
	if err := tUC.repo.GetBySlugOrId(th); err != nil {
		//logrus.Warn("thread doesn't exist")
		return tools.ThreadNotExist
	}
	v.Thread = th.Id
	if err := tUC.repo.InsertIntoVotes(v); err != nil {
		//logrus.Warn("user doesn't exist")
		return tools.UserNotExist
	}
	_ = tUC.repo.GetVotes(th, v)
	return nil
}

func (tUC *Usecase) UpdateThread(th *models.Thread) error {
	err := tUC.repo.Update(th)
	if err != nil {
		return tools.ThreadNotExist
	}
	return nil
}

func (tUC *Usecase) GetThreadPosts(th *models.Thread, desc, sort, limit, since string) ([]models.Post, error) {
	if err := tUC.repo.GetBySlugOrId(th); err != nil {
		return nil, tools.ThreadNotExist
	}
	posts, err := tUC.repo.GetPosts(th, desc, sort, limit, since)
	if err != nil {
		return nil, err
	}
	return posts, nil
}
