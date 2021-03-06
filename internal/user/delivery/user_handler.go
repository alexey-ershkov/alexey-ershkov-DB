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

	router.GET("/api/service/status", uh.GetStatus())
	router.POST("/api/service/clear", uh.DeleteAll())
	router.GET("/api/user/:nickname/profile", uh.GetUserHandler())
	router.POST("/api/user/:nickname/profile", uh.UpdateUserHandler())
	router.POST("/api/user/:nickname/create", uh.AddUserHandler())

	return uh
}

func (uh *UserHandler) AddUserHandler() echo.HandlerFunc {
	return func(c echo.Context) error {
		/*logrus.WithFields(logrus.Fields{
			"method": c.Request().Method,
		}).Info(c.Request().URL)*/
		resp := &models.User{}
		resp.Nickname = c.Param("nickname")
		err := c.Bind(resp)
		tools.HandleError(err)
		if users, err := uh.ucase.CreateUser(resp); err != nil {
			switch err {
			case tools.UserExist:
				err = c.JSON(http.StatusConflict, users)
				tools.HandleError(err)
			default:
				logrus.Error(err)
				return err
			}
		} else {
			err = c.JSON(http.StatusCreated, resp)
			tools.HandleError(err)
		}
		return nil
	}
}

func (uh *UserHandler) GetUserHandler() echo.HandlerFunc {
	return func(c echo.Context) error {
		/*logrus.WithFields(logrus.Fields{
			"method": c.Request().Method,
		}).Info(c.Request().URL)*/
		resp := &models.User{}
		resp.Nickname = c.Param("nickname")
		if err := uh.ucase.GetUser(resp); err != nil {
			switch err {
			case tools.UserNotExist:
				msg := &tools.Message{
					Message: "user not found",
				}
				err = c.JSON(http.StatusNotFound, msg)
				tools.HandleError(err)
			default:
				logrus.Error(err)
				return err
			}
		} else {
			err := c.JSON(http.StatusOK, resp)
			tools.HandleError(err)
		}
		return nil
	}
}

func (uh *UserHandler) UpdateUserHandler() echo.HandlerFunc {
	return func(c echo.Context) error {
		/*logrus.WithFields(logrus.Fields{
			"method": c.Request().Method,
		}).Info(c.Request().URL)*/
		u := &models.User{}
		u.Nickname = c.Param("nickname")
		if err := c.Bind(u); err != nil {
			tools.HandleError(err)
		}
		if err := uh.ucase.UpdateUser(u); err != nil {
			switch err {
			case tools.UserNotExist:
				err := c.JSON(
					http.StatusNotFound,
					tools.Message{
						Message: "user doesn't exist",
					})
				tools.HandleError(err)
			case tools.UserNotUpdated:
				err := c.JSON(
					http.StatusConflict,
					tools.Message{
						Message: "conflict while updating",
					})
				tools.HandleError(err)
			default:
				logrus.Error(err)
				return err
			}
		} else {
			err := c.JSON(http.StatusOK, u)
			tools.HandleError(err)
		}
		return nil
	}
}

func (uh *UserHandler) DeleteAll() echo.HandlerFunc {
	return func(c echo.Context) error {
		/*logrus.WithFields(logrus.Fields{
			"method": c.Request().Method,
		}).Info(c.Request().URL)*/
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
