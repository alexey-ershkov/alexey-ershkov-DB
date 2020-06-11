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

	tx, err := tUC.repo.CreateTx()
	if err != nil {
		return err
	}

	if err := tUC.repo.InsertInto(tx, th); err != nil {
		if err := tUC.repo.GetBySlugOrId(tx, th); err != nil {
			return tools.UserNotExist
		}

		err = tUC.repo.CommitTx(tx)
		if err != nil {
			return err
		}

		return tools.ThreadExist
	}
	if err := tUC.repo.GetCreated(tx, th); err != nil {

		err = tUC.repo.CommitTx(tx)
		if err != nil {
			return err
		}

		return err
	}

	err = tUC.repo.CommitTx(tx)
	if err != nil {
		return err
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
