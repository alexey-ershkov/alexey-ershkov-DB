package main

import (
	"alexey-ershkov/alexey-ershkov-DB.git/internal/user/delivery"
	"alexey-ershkov/alexey-ershkov-DB.git/internal/user/repository"
	"alexey-ershkov/alexey-ershkov-DB.git/internal/user/usecase"
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

	//connStr := "user=farcoad password=postgres dbname=forum sslmode=disable"
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

	rep := repository.NewUserRepo(dbConn)
	uc := usecase.NewUserUsecase(rep)
	delivery.NewUserHandler(uc, server)

	logrus.Fatal(server.Start(":5000"))
}
