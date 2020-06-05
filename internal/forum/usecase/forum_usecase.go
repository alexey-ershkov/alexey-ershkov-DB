package usecase

import (
	"alexey-ershkov/alexey-ershkov-DB.git/internal/forum"
	"alexey-ershkov/alexey-ershkov-DB.git/internal/models"
	"alexey-ershkov/alexey-ershkov-DB.git/internal/tools"
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
	tx, err := uc.repo.CreateTx()
	if err != nil {
		return err
	}

	err = uc.repo.InsertInto(tx, f)
	if err != nil {
		if err := uc.repo.GetBySlug(tx, f); err != nil {

			return tools.UserNotExist
		} else {
			err = uc.repo.CommitTx(tx)
			if err != nil {
				return err
			}

			return tools.ForumExist
		}
	}
	if err := uc.repo.GetBySlug(tx, f); err != nil {
		err = uc.repo.CommitTx(tx)
		if err != nil {
			return err
		}

		return err
	}

	err = uc.repo.CommitTx(tx)
	if err != nil {
		return err
	}

	return nil
}

func (uc *Usecase) GetForum(f *models.Forum) error {
	tx, err := uc.repo.CreateTx()
	if err != nil {
		return err
	}

	err = uc.repo.GetBySlug(tx, f)
	if err != nil {
		err = uc.repo.CommitTx(tx)
		if err != nil {
			return err
		}

		return tools.ForumNotExist
	}
	err = uc.repo.CommitTx(tx)
	if err != nil {
		return err
	}

	return nil
}

func (uc *Usecase) GetForumThreads(f *models.Forum, desc, limit, since string) ([]models.Thread, error) {
	tx, err := uc.repo.CreateTx()
	if err != nil {
		return nil, err
	}

	err = uc.repo.GetBySlug(tx, f)
	if err != nil {
		err = uc.repo.CommitTx(tx)
		if err != nil {
			return nil, err
		}

		return nil, tools.ForumNotExist
	}

	ths, err := uc.repo.GetThreads(tx, f, desc, limit, since)
	if err != nil {

		err = uc.repo.CommitTx(tx)
		if err != nil {
			return nil, err
		}

		return nil, err
	}

	err = uc.repo.CommitTx(tx)
	if err != nil {
		return nil, err
	}

	return ths, nil
}

func (uc *Usecase) GetForumUsers(f *models.Forum, desc, limit, since string) ([]models.User, error) {
	tx, err := uc.repo.CreateTx()
	if err != nil {
		return nil, err
	}

	err = uc.repo.GetBySlug(tx, f)
	if err != nil {
		err = uc.repo.CommitTx(tx)
		if err != nil {
			return nil, err
		}

		return nil, tools.ForumNotExist
	}
	usr, err := uc.repo.GetUsers(tx, f, desc, limit, since)
	if err != nil {
		err = uc.repo.CommitTx(tx)
		if err != nil {
			return nil, err
		}

		return nil, err
	}

	err = uc.repo.CommitTx(tx)
	if err != nil {
		return nil, err
	}

	return usr, nil
}
