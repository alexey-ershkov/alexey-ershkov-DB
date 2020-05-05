package main

import (
	fHandler "alexey-ershkov/alexey-ershkov-DB.git/internal/forum/delivery"
	fRepo "alexey-ershkov/alexey-ershkov-DB.git/internal/forum/repository"
	fUUcase "alexey-ershkov/alexey-ershkov-DB.git/internal/forum/usecase"
	thHandler "alexey-ershkov/alexey-ershkov-DB.git/internal/thread/delivery"
	thRepo "alexey-ershkov/alexey-ershkov-DB.git/internal/thread/repository"
	thUcase "alexey-ershkov/alexey-ershkov-DB.git/internal/thread/usecase"
	uHanler "alexey-ershkov/alexey-ershkov-DB.git/internal/user/delivery"
	uRepo "alexey-ershkov/alexey-ershkov-DB.git/internal/user/repository"
	uUcase "alexey-ershkov/alexey-ershkov-DB.git/internal/user/usecase"
	"github.com/jackc/pgx"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"os"
)

func init() {
	logrus.SetOutput(os.Stdout)
	logrus.SetLevel(logrus.DebugLevel)
}

func main() {
	server := echo.New()

	dbConf := pgx.ConnConfig{
		User:                 "farcoad",
		Database:             "forum",
		Password:             "postgres",
		PreferSimpleProtocol: true,
	}
	dbConn, err := pgx.Connect(dbConf)
	if err != nil {
		logrus.Fatal(err)
	}

	uRep := uRepo.NewUserRepo(dbConn)
	fRep := fRepo.NewForumRepository(dbConn)
	thRep := thRepo.NewThreadRepository(dbConn)

	uUC := uUcase.NewUserUsecase(uRep)
	fUC := fUUcase.NewForumUsecase(fRep)
	thUC := thUcase.NewThreadUsecase(thRep)

	uHanler.NewUserHandler(uUC, server)
	fHandler.NewForumHandler(server, fUC)
	thHandler.NewThreadHandler(thUC, server)

	logrus.Fatal(server.Start(":5000"))
}
