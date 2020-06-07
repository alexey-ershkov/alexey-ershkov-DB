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
	for _, val := range p {
		val.Thread = th.Id
		val.Forum = th.Forum
		u := &models.User{}
		u.Nickname = val.Author
		val.Author = u.Nickname
		if val.Parent != 0 {
			sp := &models.Post{}
			sp.Id = val.Parent
			if err = pUC.pRepo.GetById(tx, sp); err != nil {
				return tools.ParentNotExist
			}
			if sp.Thread != val.Thread {
				return tools.ParentNotExist
			}
			val.Path = sp.Path
		}
	}
	if err = pUC.pRepo.InsertInto(tx, p); err != nil {
		return tools.UserNotExist
	}

	return nil
}

func (pUC *PostUsecase) GetPost(p *models.Post) error {
	tx, err := pUC.pRepo.CreateTx()
	if err != nil {
		return err
	}

	if err := pUC.pRepo.GetById(tx, p); err != nil {
		//logrus.Warn("post not exist")

		err = pUC.pRepo.CommitTx(tx)
		if err != nil {
			return err
		}

		return tools.PostNotExist
	}

	err = pUC.pRepo.CommitTx(tx)
	if err != nil {
		return err
	}

	return nil
}

func (pUC *PostUsecase) UpdatePost(p *models.Post) error {

	tx, err := pUC.pRepo.CreateTx()
	if err != nil {
		return err
	}

	message := p.Message
	if err := pUC.pRepo.GetById(tx, p); err != nil {
		return tools.PostNotExist
	}
	if message != "" && message != p.Message {
		p.Message = message
		if err := pUC.pRepo.Update(tx, p); err != nil {

			err = pUC.pRepo.CommitTx(tx)
			if err != nil {
				return err
			}

			return tools.PostNotExist
		}
	}

	err = pUC.pRepo.CommitTx(tx)
	if err != nil {
		return err
	}

	return nil
}
