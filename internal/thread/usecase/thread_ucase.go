package usecase

import (
	"alexey-ershkov/alexey-ershkov-DB.git/internal/forum"
	"alexey-ershkov/alexey-ershkov-DB.git/internal/models"
	"alexey-ershkov/alexey-ershkov-DB.git/internal/thread"
	"alexey-ershkov/alexey-ershkov-DB.git/internal/tools"
)

type Usecase struct {
	repo thread.Repository
	fRepo forum.Repository
}

func NewThreadUsecase(r thread.Repository, fr forum.Repository) thread.Usecase {
	return &Usecase{
		repo: r,
		fRepo: fr,
	}
}

func (tUC *Usecase) CreateThread(th *models.Thread) error {
	tx, err := tUC.repo.CreateTx()
	defer func() {
		if err == nil {
			_ = tx.Commit()
		} else {
			_ = tx.Rollback()
		}
	}()
	if err != nil {
		return err
	}

	f := &models.Forum{}
	f.Slug = th.Forum
	err = tUC.fRepo.GetBySlug(tx, f)
	if err != nil {
		return tools.UserNotExist
	}

	th.Forum = f.Slug

	err = tUC.repo.InsertInto(tx, th);
	if err != nil {
		err = tUC.repo.GetBySlugOrId(tx, th);
		if err != nil {
			return tools.UserNotExist
		}

		return tools.ThreadExist
	}

	return nil
}

func (tUC *Usecase) GetThreadInfo(th *models.Thread) error {

	tx, err := tUC.repo.CreateTx()
	defer func() {
		if err == nil {
			_ = tx.Commit()
		} else {
			_ = tx.Rollback()
		}
	}()
	if err != nil {
		return err
	}

	err = tUC.repo.GetBySlugOrId(tx, th)
	if err != nil {
		return tools.ThreadNotExist
	}

	return nil
}

func (tUC *Usecase) CreateVote(th *models.Thread, v *models.Vote) error {
	tx, err := tUC.repo.CreateTx()
	defer func() {
		if err == nil {
			_ = tx.Commit()
		} else {
			_ = tx.Rollback()
		}
	}()
	if err != nil {
		return err
	}

	err = tUC.repo.GetBySlugOrId(tx, th)
	if err != nil {
		return tools.ThreadNotExist
	}
	v.Thread = th.Id
	err = tUC.repo.InsertIntoVotes(tx, th,v)
	if err != nil {
		return tools.UserNotExist
	}

	return nil
}

func (tUC *Usecase) UpdateThread(th *models.Thread) error {
	tx, err := tUC.repo.CreateTx()
	defer func() {
		if err == nil {
			_ = tx.Commit()
		} else {
			_ = tx.Rollback()
		}
	}()
	if err != nil {
		return err
	}

	err = tUC.repo.Update(tx, th)
	if err != nil {
		return tools.ThreadNotExist
	}

	return nil
}

func (tUC *Usecase) GetThreadPosts(th *models.Thread, desc, sort, limit, since string) ([]models.Post, error) {
	tx, err := tUC.repo.CreateTx()
	defer func() {
		if err == nil {
			_ = tx.Commit()
		} else {
			_ = tx.Rollback()
		}
	}()
	if err != nil {
		return nil, err
	}

	err = tUC.repo.GetBySlugOrId(tx, th)
	if err != nil {
		return nil, tools.ThreadNotExist
	}
	posts, e := tUC.repo.GetPosts(tx, th, desc, sort, limit, since)
	err = e
	if err != nil {
		return nil, err
	}
	return posts, nil
}
