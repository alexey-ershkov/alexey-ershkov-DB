package usecase

import (
	"alexey-ershkov/alexey-ershkov-DB.git/internal/forum"
	"alexey-ershkov/alexey-ershkov-DB.git/internal/models"
	"alexey-ershkov/alexey-ershkov-DB.git/internal/tools"
	"alexey-ershkov/alexey-ershkov-DB.git/internal/user"
)

type Usecase struct {
	repo forum.Repository
	uRep user.Repository
}

func NewForumUsecase(r forum.Repository, ur user.Repository) forum.Usecase {
	return &Usecase{
		repo: r,
		uRep: ur,
	}
}

func (uc *Usecase) CreateForum(f *models.Forum) error {
	tx, err := uc.repo.CreateTx()
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

	u := &models.User{}
	u.Nickname = f.User
	err = uc.uRep.GetByNickname(tx, u)
	if err != nil {

		return tools.UserNotExist
	}

	f.User = u.Nickname

	err = uc.repo.GetBySlug(tx, f)
	if err == nil {
		return tools.ForumExist
	}

	err = uc.repo.InsertInto(tx, f)
	if err != nil {
		return err
	}

	return nil
}

func (uc *Usecase) GetForum(f *models.Forum) error {
	tx, err := uc.repo.CreateTx()
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

	err = uc.repo.GetBySlug(tx, f)
	if err != nil {
		return tools.ForumNotExist
	}

	return nil
}

func (uc *Usecase) GetForumThreads(f *models.Forum, desc, limit, since string) ([]models.Thread, error) {
	tx, err := uc.repo.CreateTx()
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

	err = uc.repo.GetBySlug(tx, f)
	if err != nil {
		return nil, tools.ForumNotExist
	}

	ths, err := uc.repo.GetThreads(tx, f, desc, limit, since)
	if err != nil {
		return nil, err
	}

	return ths, nil
}

func (uc *Usecase) GetForumUsers(f *models.Forum, desc, limit, since string) ([]models.User, error) {
	tx, err := uc.repo.CreateTx()
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

	err = uc.repo.GetBySlug(tx, f)
	if err != nil {
		return nil, tools.ForumNotExist
	}
	usr, err := uc.repo.GetUsers(tx, f, desc, limit, since)
	if err != nil {
		return nil, err
	}

	return usr, nil
}
