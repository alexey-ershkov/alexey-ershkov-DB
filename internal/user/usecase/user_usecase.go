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
	err := uc.Repo.InsertInto(u)
	if err != nil {
		//logrus.Warn("User already exist")
		users, err := uc.Repo.GetByNicknameOrEmail(u)
		tools.HandleError(err)
		return users, tools.UserExist
	}
	return nil, nil
}

func (uc *UserUsecase) GetUser(u *models.User) error {
	err := uc.Repo.GetByNickname(u)
	if err != nil {
		return tools.UserNotExist
	}
	return nil
}

func (uc *UserUsecase) UpdateUser(u *models.User) error {
	uInfo := *u
	if err := uc.Repo.GetByNickname(&uInfo); err != nil {
		return tools.UserNotExist
	}
	if u.Email == "" {
		u.Email = uInfo.Email
	}
	if err := uc.Repo.Update(u); err != nil {
		return tools.UserNotUpdated
	}
	if u.About == "" {
		u.About = uInfo.About
	}
	if u.Fullname == "" {
		u.Fullname = uInfo.Fullname
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
	err := uc.Repo.GetStatus(s)
	if err != nil {
		return err
	}
	return nil
}
