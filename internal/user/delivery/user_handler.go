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

	router.GET("/service/status", uh.GetStatus())
	router.POST("/service/clear", uh.DeleteAll())
	router.GET("/user/:nickname/profile", uh.GetUserHandler())
	router.POST("/user/:nickname/profile", uh.UpdateUserHandler())
	router.POST("/user/:nickname/create", uh.AddUserHandler())

	return uh
}

func (uh *UserHandler) AddUserHandler() echo.HandlerFunc {
	return func(c echo.Context) error {
		logrus.WithFields(logrus.Fields{
			"method": c.Request().Method,
		}).Info(c.Request().URL)
		resp := &models.User{}
		resp.Nickname = c.Param("nickname")
		err := c.Bind(resp)
		tools.HandleError(err)
		if users, err := uh.ucase.CreateUser(resp); err != nil {
			if err == tools.UserExist {
				err = c.JSON(http.StatusConflict, users)
				tools.HandleError(err)
				return nil
			}
			logrus.Error(err)
			return err
		}

		err = c.JSON(http.StatusCreated, resp)
		tools.HandleError(err)
		return nil
	}
}

func (uh *UserHandler) GetUserHandler() echo.HandlerFunc {
	return func(c echo.Context) error {
		logrus.WithFields(logrus.Fields{
			"method": c.Request().Method,
		}).Info(c.Request().URL)
		resp := &models.User{}
		resp.Nickname = c.Param("nickname")
		if err := uh.ucase.GetUser(resp); err != nil {
			msg := &tools.Message{
				Message: "user not found",
			}
			err = c.JSON(http.StatusNotFound, msg)
			tools.HandleError(err)
			return nil
		}
		err := c.JSON(http.StatusOK, resp)
		tools.HandleError(err)
		return nil
	}
}

func (uh *UserHandler) UpdateUserHandler() echo.HandlerFunc {
	return func(c echo.Context) error {
		logrus.WithFields(logrus.Fields{
			"method": c.Request().Method,
		}).Info(c.Request().URL)
		u := &models.User{}
		u.Nickname = c.Param("nickname")
		if err := uh.ucase.GetUser(u); err != nil {
			if err := c.JSON(http.StatusNotFound, tools.Message{
				"user doesn't exist",
			}); err != nil {
				tools.HandleError(err)
			}
			return nil
		}
		if err := c.Bind(u); err != nil {
			tools.HandleError(err)
		}
		if err := uh.ucase.UpdateUser(u); err != nil {
			if err := c.JSON(http.StatusConflict, tools.Message{
				Message: "conflict while updating",
			}); err != nil {
				tools.HandleError(err)
			}
			return nil
		}
		if err := c.JSON(http.StatusOK, u); err != nil {
			tools.HandleError(err)
		}
		return nil
	}
}

func (uh *UserHandler) DeleteAll() echo.HandlerFunc {
	return func(c echo.Context) error {
		logrus.WithFields(logrus.Fields{
			"method": c.Request().Method,
		}).Info(c.Request().URL)
		err := uh.ucase.DeleteAll()
		tools.HandleError(err)
		err = c.JSON(http.StatusOK, tools.Message{
			Message: "all info deleted",
		})
		tools.HandleError(err)
		return nil
	}
}

func (uh *UserHandler) GetStatus() echo.HandlerFunc {
	return func(c echo.Context) error {
		s := &models.Status{}
		err := uh.ucase.GetStatus(s)
		tools.HandleError(err)
		err = c.JSON(http.StatusOK, s)
		tools.HandleError(err)
		return nil
	}
}
