package delivery

import (
	"alexey-ershkov/alexey-ershkov-DB.git/internal/models"
	"alexey-ershkov/alexey-ershkov-DB.git/internal/tools"
	"alexey-ershkov/alexey-ershkov-DB.git/internal/user"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"net/http"
)

type UserHandler struct {
	ucase user.Usecase
}

func NewUserHandler(uc user.Usecase, router *echo.Echo) *UserHandler {
	uh := &UserHandler{
		ucase: uc,
	}

	router.GET("/user/:nickname/profile", uh.GetUser())
	router.POST("/user/:nickname/create", uh.AddUser())

	return uh
}

func (uh *UserHandler) AddUser() echo.HandlerFunc {
	return func(c echo.Context) error {
		logrus.Info(c.Request().URL)
		resp := &models.User{}
		resp.Nickname = c.Param("nickname")
		err := c.Bind(resp)
		handleError(err)

		if users, err := uh.ucase.CreateUser(resp); err != nil {
			if err == tools.UserExist {
				err = c.JSON(http.StatusConflict, users)
				handleError(err)
				return nil
			}
			logrus.Error(err)
			return err
		}

		err = c.JSON(http.StatusCreated, resp)
		handleError(err)
		return nil
	}
}

func (uh *UserHandler) GetUser() echo.HandlerFunc {
	return func(c echo.Context) error {
		logrus.Info(c.Request().URL)
		resp := &models.User{}
		resp.Nickname = c.Param("nickname")
		if err := uh.ucase.GetUser(resp); err != nil {
			msg := &tools.Message{
				Message: "user not found",
			}
			err = c.JSON(http.StatusNotFound, msg)
			handleError(err)
			return nil
		}
		err := c.JSON(http.StatusOK, resp)
		handleError(err)
		return nil
	}
}

func handleError(e error) {
	if e != nil {
		logrus.Error(e)
	}
}
