package delivery

import (
	"alexey-ershkov/alexey-ershkov-DB.git/internal/models"
	"alexey-ershkov/alexey-ershkov-DB.git/internal/thread"
	"alexey-ershkov/alexey-ershkov-DB.git/internal/tools"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"net/http"
)

type ThreadHandler struct {
	thUC thread.Usecase
}

func NewThreadHandler(uc thread.Usecase, router *echo.Echo) {
	thH := &ThreadHandler{
		thUC: uc,
	}

	router.POST("forum/:forum/create", thH.CreateThread())
	router.GET("thread/:slug/details", thH.GetThreadInfo())
	router.POST("thread/:slug/details", thH.UpdateThread())
	router.POST("thread/:slug/vote", thH.CreateVote())
}

func (thH *ThreadHandler) CreateThread() echo.HandlerFunc {
	return func(c echo.Context) error {
		logrus.WithFields(logrus.Fields{
			"method": c.Request().Method,
		}).Info(c.Request().URL)
		th := &models.Thread{}
		th.Forum = c.Param("forum")
		err := c.Bind(th)
		tools.HandleError(err)
		if err := thH.thUC.CreateThread(th); err != nil {
			if err == tools.ThreadExist {
				logrus.Warn("thread already exists")
				err := c.JSON(http.StatusConflict, th)
				tools.HandleError(err)
				return nil
			}
			if err == tools.UserNotExist {
				logrus.Warn("user not exists")
				err = c.JSON(http.StatusNotFound, tools.Message{
					Message: "user not exist",
				})
				tools.HandleError(err)
				return nil
			}
			tools.HandleError(err)
			return err
		}
		err = c.JSON(http.StatusCreated, th)
		tools.HandleError(err)
		return nil
	}
}

func (thH *ThreadHandler) GetThreadInfo() echo.HandlerFunc {
	return func(c echo.Context) error {
		logrus.WithFields(logrus.Fields{
			"method": c.Request().Method,
		}).Info(c.Request().URL)
		th := &models.Thread{}
		th.Slug = c.Param("slug")
		err := c.Bind(th)
		tools.HandleError(err)
		err = thH.thUC.GetThreadInfo(th)
		if err == tools.ThreadNotExist {
			err := c.JSON(http.StatusNotFound, tools.Message{
				Message: "thread not found",
			})
			tools.HandleError(err)
			return nil
		}
		err = c.JSON(http.StatusOK, th)
		tools.HandleError(err)
		return nil
	}
}

func (thH *ThreadHandler) CreateVote() echo.HandlerFunc {
	return func(c echo.Context) error {
		logrus.WithFields(logrus.Fields{
			"method": c.Request().Method,
		}).Info(c.Request().URL)
		th := &models.Thread{}
		vote := &models.Vote{}
		th.Slug = c.Param("slug")
		err := c.Bind(vote)
		tools.HandleError(err)
		err = thH.thUC.CreateVote(th, vote)
		if err == tools.UserNotExist {
			e := c.JSON(http.StatusNotFound, tools.Message{
				Message: "user not found",
			})
			tools.HandleError(e)
			return nil
		}
		if err == tools.ThreadNotExist {
			e := c.JSON(http.StatusNotFound, tools.Message{
				Message: "thread not found",
			})
			tools.HandleError(e)
			return nil
		}
		err = c.JSON(http.StatusOK, th)
		tools.HandleError(err)
		return nil
	}
}

func (thH *ThreadHandler) UpdateThread() echo.HandlerFunc {
	return func(c echo.Context) error {
		logrus.WithFields(logrus.Fields{
			"method": c.Request().Method,
		}).Info(c.Request().URL)
		th := &models.Thread{}
		th.Slug = c.Param("slug")
		err := c.Bind(th)
		tools.HandleError(err)
		err = thH.thUC.UpdateThread(th)
		if err == tools.ThreadNotExist {
			err := c.JSON(http.StatusNotFound, tools.Message{
				Message: "thread not found",
			})
			tools.HandleError(err)
			return nil
		}
		tools.HandleError(err)
		err = c.JSON(http.StatusOK, th)
		tools.HandleError(err)
		return nil
	}
}
