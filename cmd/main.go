package main

import (
	fHandler "alexey-ershkov/alexey-ershkov-DB.git/internal/forum/delivery"
	fRepo "alexey-ershkov/alexey-ershkov-DB.git/internal/forum/repository"
	fUUcase "alexey-ershkov/alexey-ershkov-DB.git/internal/forum/usecase"
	"time"

	thHandler "alexey-ershkov/alexey-ershkov-DB.git/internal/thread/delivery"
	thRepo "alexey-ershkov/alexey-ershkov-DB.git/internal/thread/repository"
	thUcase "alexey-ershkov/alexey-ershkov-DB.git/internal/thread/usecase"

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

func main() {
	server := echo.New()

	dbConf := pgx.ConnConfig{
		User:                 "farcoad",
		Database:             "forum",
		Password:             "postgres",
		PreferSimpleProtocol: false,
	}

	dbPoolConf := pgx.ConnPoolConfig{
		ConnConfig:     dbConf,
		MaxConnections: 100,
		AfterConnect:   nil,
		AcquireTimeout: 0,
	}

	dbConn, err := pgx.NewConnPool(dbPoolConf)
	if err != nil {
		logrus.Fatal(err)
	}

	StatsLog(dbConn)

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

func StatsLog(conn *pgx.ConnPool) {
	ticker := time.NewTicker(time.Second * 5)
	go func() {
		for range ticker.C {
			stats := conn.Stat()
			logrus.WithFields(logrus.Fields{
				"Max Conn":       stats.MaxConnections,
				"Current Conn":   stats.CurrentConnections,
				"Avaliable Conn": stats.AvailableConnections,
			}).Info()
		}
	}()
}
