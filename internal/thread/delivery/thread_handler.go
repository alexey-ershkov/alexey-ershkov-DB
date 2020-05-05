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
}

func (thH *ThreadHandler) CreateThread() echo.HandlerFunc {
	return func(c echo.Context) error {
		logrus.Info(c.Request().Method, "   ", c.Request().URL)
		th := &models.Thread{}
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
