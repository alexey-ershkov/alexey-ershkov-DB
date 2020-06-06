package usecase

import (
	"alexey-ershkov/alexey-ershkov-DB.git/internal/models"
	"alexey-ershkov/alexey-ershkov-DB.git/internal/tools"
	"alexey-ershkov/alexey-ershkov-DB.git/internal/user"
)

type UserUsecase struct {
	Repo user.Repository
}

func NewUserUsecase(r user.Repository) user.Usecase {
	return &UserUsecase{
		Repo: r,
	}
}

func (uc *UserUsecase) CreateUser(u *models.User) ([]models.User, error) {
	tx, err := uc.Repo.CreateTx()
	if err != nil {
		return nil, err
	}

	err = uc.Repo.InsertInto(tx, u)
	if err != nil {
		//logrus.Warn("User already exist")
		users, err := uc.Repo.GetByNicknameOrEmail(tx, u)
		tools.HandleError(err)
		return users, tools.UserExist
	}

	err = uc.Repo.CommitTx(tx)
	if err != nil {
		return nil, err
	}
	return nil, nil
}

func (uc *UserUsecase) GetUser(u *models.User) error {
	tx, err := uc.Repo.CreateTx()
	if err != nil {
		return err
	}

	err = uc.Repo.GetByNickname(tx, u)
	if err != nil {
		err = uc.Repo.CommitTx(tx)
		if err != nil {
			return err
		}
		return tools.UserNotExist
	}

	err = uc.Repo.CommitTx(tx)
	if err != nil {
		return err
	}
	return nil
}

func (uc *UserUsecase) UpdateUser(u *models.User) error {
	tx, err := uc.Repo.CreateTx()
	if err != nil {
		return err
	}

	uInfo := *u
	if err := uc.Repo.GetByNickname(tx, &uInfo); err != nil {
		return tools.UserNotExist
	}
	if u.Email == "" {
		u.Email = uInfo.Email
	}
	if u.About == "" {
		u.About = uInfo.About
	}
	if u.Fullname == "" {
		u.Fullname = uInfo.Fullname
	}
	if err := uc.Repo.Update(tx, u); err != nil {
		return tools.UserNotUpdated
	}

	err = uc.Repo.CommitTx(tx)
	if err != nil {
		return err
	}
	return nil
}

func (uc *UserUsecase) DeleteAll() error {
	err := uc.Repo.DeleteAll()
	if err != nil {
		return err
	}
	return nil
}

func (uc *UserUsecase) GetStatus(s *models.Status) error {
	tx, err := uc.Repo.CreateTx()
	if err != nil {
		return err
	}

	err = uc.Repo.GetStatus(tx, s)
	if err != nil {
		return err
	}

	err = uc.Repo.CommitTx(tx)
	if err != nil {
		return err
	}
	return nil
}
