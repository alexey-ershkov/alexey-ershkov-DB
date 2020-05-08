package delivery

import (
	"alexey-ershkov/alexey-ershkov-DB.git/internal/models"
	"alexey-ershkov/alexey-ershkov-DB.git/internal/post"
	"alexey-ershkov/alexey-ershkov-DB.git/internal/tools"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"net/http"
	"strconv"
)

type PostHandler struct {
	pUC post.Usecase
}

func NewPostHandler(router *echo.Echo, pUC post.Usecase) {
	ph := &PostHandler{
		pUC: pUC,
	}

	router.POST("/thread/:slug/create", ph.CreatePosts())
	router.GET("/post/:id", ph.GetPost())
}

func (ph *PostHandler) CreatePosts() echo.HandlerFunc {
	return func(c echo.Context) error {
		logrus.WithFields(logrus.Fields{
			"method": c.Request().Method,
		}).Info(c.Request().URL)
		p := make([]*models.Post, 0)
		err := c.Bind(&p)
		tools.HandleError(err)
		th := models.Thread{}
		th.Slug = c.Param("slug")
		err = ph.pUC.CreatePosts(p, &th)
		if err != nil {
			if err == tools.ParentNotExist {
				e := c.JSON(http.StatusConflict, tools.Message{
					Message: "parent conflict",
				})
				tools.HandleError(e)
				return nil
			}
			if err == tools.ThreadNotExist || err == tools.UserNotExist {
				e := c.JSON(http.StatusNotFound, tools.Message{
					Message: "user or thread not found",
				})
				tools.HandleError(e)
				return nil
			}
			e := c.JSON(http.StatusInternalServerError, tools.Message{
				Message: err.Error(),
			})
			tools.HandleError(e)
			return nil
		}
		err = c.JSON(http.StatusCreated, p)
		tools.HandleError(err)
		return nil
	}
}

func (ph *PostHandler) GetPost() echo.HandlerFunc {
	return func(c echo.Context) error {
		logrus.WithFields(logrus.Fields{
			"method": c.Request().Method,
		}).Info(c.Request().URL)
		p := models.Post{}
		var err error
		p.Id, err = strconv.ParseInt(c.Param("id"), 10, 64)
		if err != nil {
			logrus.Error("Cannot parse id")
		}
		if err := ph.pUC.GetPost(&p); err != nil {
			err := c.JSON(http.StatusNotFound, tools.Message{
				Message: err.Error(),
			})
			tools.HandleError(err)
			return nil
		}
		err = c.JSON(http.StatusOK, p)
		tools.HandleError(err)
		return nil
	}
}
