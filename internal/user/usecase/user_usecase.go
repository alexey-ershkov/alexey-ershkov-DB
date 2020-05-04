package usecase

import (
	"alexey-ershkov/alexey-ershkov-DB.git/internal/models"
	"alexey-ershkov/alexey-ershkov-DB.git/internal/tools"
	"alexey-ershkov/alexey-ershkov-DB.git/internal/user"
	"github.com/sirupsen/logrus"
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
	err := uc.Repo.InsertInto(u)
	if err != nil {
		logrus.Warn("User already exist")
		users, err := uc.Repo.GetByNicknameOrEmail(u)
		if err != nil {
			logrus.Error(err)
		}
		return users, tools.UserExist
	}
	return nil, nil
}

func (uc *UserUsecase) GetUser(u *models.User) error {
	err := uc.Repo.GetByNickname(u)
	if err != nil {
		return err
	}
	return nil
}

func (uc *UserUsecase) UpdateUser(u *models.User) error {
	if err := uc.Repo.Update(u); err != nil {
		return err
	}
	return nil
}
