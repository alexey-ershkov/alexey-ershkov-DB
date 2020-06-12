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
		err = tUC.repo.GetBySlug(tx, th);
		if err != nil {
			return tools.UserNotExist
		}

		return tools.ThreadExist
	}

	return nil
}

func (tUC *Usecase) GetThreadInfo(th *models.Thread) error {

	tx, err := tUC.repo.CreateTx()
	if err != nil {
		return err
	}

	if err := tUC.repo.GetBySlugOrId(tx, th); err != nil {
		//logrus.Warn("thread doesn't exist")

		err = tUC.repo.CommitTx(tx)
		if err != nil {
			return err
		}

		return tools.ThreadNotExist
	}

	err = tUC.repo.CommitTx(tx)
	if err != nil {
		return err
	}

	return nil
}

func (tUC *Usecase) CreateVote(th *models.Thread, v *models.Vote) error {
	tx, err := tUC.repo.CreateTx()
	if err != nil {
		return err
	}

	if err := tUC.repo.GetBySlugOrId(tx, th); err != nil {
		//logrus.Warn("thread doesn't exist")

		err = tUC.repo.CommitTx(tx)
		if err != nil {
			return err
		}

		return tools.ThreadNotExist
	}
	v.Thread = th.Id
	if err := tUC.repo.InsertIntoVotes(tx, v); err != nil {
		//logrus.Warn("user doesn't exist")

		return tools.UserNotExist
	}
	_ = tUC.repo.GetVotes(tx, th, v)

	err = tUC.repo.CommitTx(tx)
	if err != nil {
		return err
	}

	return nil
}

func (tUC *Usecase) UpdateThread(th *models.Thread) error {
	tx, err := tUC.repo.CreateTx()
	if err != nil {
		return err
	}

	err = tUC.repo.Update(tx, th)
	if err != nil {

		err = tUC.repo.CommitTx(tx)
		if err != nil {
			return err
		}

		return tools.ThreadNotExist
	}

	err = tUC.repo.CommitTx(tx)
	if err != nil {
		return err
	}

	return nil
}

func (tUC *Usecase) GetThreadPosts(th *models.Thread, desc, sort, limit, since string) ([]models.Post, error) {
	tx, err := tUC.repo.CreateTx()
	if err != nil {
		return nil, err
	}

	if err := tUC.repo.GetBySlugOrId(tx, th); err != nil {

		err = tUC.repo.CommitTx(tx)
		if err != nil {
			return nil, err
		}

		return nil, tools.ThreadNotExist
	}
	posts, err := tUC.repo.GetPosts(tx, th, desc, sort, limit, since)
	if err != nil {

		err = tUC.repo.CommitTx(tx)
		if err != nil {
			return nil, err
		}

		return nil, err
	}

	err = tUC.repo.CommitTx(tx)
	if err != nil {
		return nil, err
	}

	return posts, nil
}
