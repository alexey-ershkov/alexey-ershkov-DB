package usecase

import (
	"alexey-ershkov/alexey-ershkov-DB.git/internal/models"
	"alexey-ershkov/alexey-ershkov-DB.git/internal/post"
	"alexey-ershkov/alexey-ershkov-DB.git/internal/thread"
	"alexey-ershkov/alexey-ershkov-DB.git/internal/tools"
	"alexey-ershkov/alexey-ershkov-DB.git/internal/user"
)

type PostUsecase struct {
	pRepo post.Repository
	tRepo thread.Repository
	uRepo user.Repository
}

func NewPostUsecase(pr post.Repository, tr thread.Repository, ur user.Repository) post.Usecase {
	return &PostUsecase{
		pRepo: pr,
		tRepo: tr,
		uRepo: ur,
	}
}

func (pUC *PostUsecase) CreatePosts(p []*models.Post, th *models.Thread) error {
	tx, err := pUC.pRepo.CreateTx()
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

	if err = pUC.tRepo.GetBySlugOrId(tx, th); err != nil {
		return tools.ThreadNotExist
	}
	if err = pUC.pRepo.InsertInto(tx, p, th); err != nil {
		if err.Error() == "ERROR: Parent post was created in another thread (SQLSTATE 00404)" {
			return tools.ParentNotExist
		} else {
			return tools.UserNotExist
		}
	}

	for iter := range p {
		pUC.tRepo.InsertIntoForumUsers(tx, p[iter].Forum, p[iter].Author)
	}

	return nil
}

func (pUC *PostUsecase) GetPost(p *models.Post) error {
	tx, err := pUC.pRepo.CreateTx()
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

	err = pUC.pRepo.GetById(tx, p)
	if err != nil {
		return tools.PostNotExist
	}

	return nil
}

func (pUC *PostUsecase) UpdatePost(p *models.Post) error {

	tx, err := pUC.pRepo.CreateTx()
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

	message := p.Message
	err = pUC.pRepo.GetById(tx, p)
	if err != nil {
		return tools.PostNotExist
	}
	if message != "" && message != p.Message {
		p.Message = message
		if err := pUC.pRepo.Update(tx, p); err != nil {
			return tools.PostNotExist
		}
	}

	return nil
}
