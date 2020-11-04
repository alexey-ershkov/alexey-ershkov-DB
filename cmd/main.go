package main

import (
	fHandler "alexey-ershkov/alexey-ershkov-DB.git/internal/forum/delivery"
	fRepo "alexey-ershkov/alexey-ershkov-DB.git/internal/forum/repository"
	fUUcase "alexey-ershkov/alexey-ershkov-DB.git/internal/forum/usecase"
	thHandler "alexey-ershkov/alexey-ershkov-DB.git/internal/thread/delivery"
	thRepo "alexey-ershkov/alexey-ershkov-DB.git/internal/thread/repository"
	thUcase "alexey-ershkov/alexey-ershkov-DB.git/internal/thread/usecase"
	"net/http"

	uHandler "alexey-ershkov/alexey-ershkov-DB.git/internal/user/delivery"
	uRepo "alexey-ershkov/alexey-ershkov-DB.git/internal/user/repository"
	uUcase "alexey-ershkov/alexey-ershkov-DB.git/internal/user/usecase"

	pHandler "alexey-ershkov/alexey-ershkov-DB.git/internal/post/delivery"
	pRepo "alexey-ershkov/alexey-ershkov-DB.git/internal/post/repository"
	pUcase "alexey-ershkov/alexey-ershkov-DB.git/internal/post/usecase"

	"github.com/jackc/pgx"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"os"
)

func init() {
	logrus.SetOutput(os.Stdout)
	logrus.SetLevel(logrus.DebugLevel)
}

func PanicMiddleWare(next echo.HandlerFunc) echo.HandlerFunc {

	return func(c echo.Context) error {
		defer func() error {
			if err := recover(); err != nil {
				logrus.WithFields(logrus.Fields{
					"ERROR": err.(error).Error(),
				}).Error()
				return c.NoContent(http.StatusInternalServerError)
			}
			return nil
		}()

		return next(c)
	}
}

func main() {
	server := echo.New()
	server.Use(PanicMiddleWare)

	dbConf := pgx.ConnConfig{
		User:                 "farcoad",
		Database:             "forum",
		Password:             "postgres",
		PreferSimpleProtocol: false,
	}



	dbConn, err := pgx.Connect(dbConf)

	if err != nil {
		logrus.Fatal(err)
	}


	uRep := uRepo.NewUserRepo(dbConn)
	fRep := fRepo.NewForumRepository(dbConn)
	thRep := thRepo.NewThreadRepository(dbConn)
	pRep := pRepo.NewPostRepository(dbConn)

	err = uRep.Prepare()
	if err != nil {
		logrus.Fatal(err)
	}
	err = fRep.Prepare()
	if err != nil {
		logrus.Fatal(err)
	}
	err = thRep.Prepare()
	if err != nil {
		logrus.Fatal(err)
	}
	err = pRep.Prepare()
	if err != nil {
		logrus.Fatal(err)
	}

	uUC := uUcase.NewUserUsecase(uRep)
	fUC := fUUcase.NewForumUsecase(fRep)
	thUC := thUcase.NewThreadUsecase(thRep)
	pUC := pUcase.NewPostUsecase(pRep, thRep, uRep)

	uHandler.NewUserHandler(uUC, server)
	fHandler.NewForumHandler(server, fUC)
	thHandler.NewThreadHandler(thUC, server)
	pHandler.NewPostHandler(server, pUC, fUC, uUC, thUC)

	logrus.Fatal(server.Start(":5000"))
}


