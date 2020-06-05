package delivery

import (
	"alexey-ershkov/alexey-ershkov-DB.git/internal/forum"
	"alexey-ershkov/alexey-ershkov-DB.git/internal/models"
	"alexey-ershkov/alexey-ershkov-DB.git/internal/tools"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"net/http"
)

type ForumHandler struct {
	uc forum.Usecase
}

func NewForumHandler(router *echo.Echo, uc forum.Usecase) {
	fh := &ForumHandler{uc: uc}

	router.POST("/api/forum/create", fh.CreateForum())
	router.GET("/api/forum/:slug/details", fh.GetForumInfo())
	router.GET("/api/forum/:slug/threads", fh.GetForumThreads())
	router.GET("/api/forum/:slug/users", fh.GetForumUsers())
}

func (fh *ForumHandler) CreateForum() echo.HandlerFunc {
	return func(c echo.Context) error {
		/*logrus.WithFields(logrus.Fields{
			"method": c.Request().Method,
		}).Info(c.Request().URL)*/
		f := &models.Forum{}
		err := c.Bind(f)
		if err != nil {
			tools.HandleError(err)
		}
		if err := fh.uc.CreateForum(f); err != nil {
			switch err {
			case tools.UserNotExist:
				err := c.JSON(http.StatusNotFound, tools.Message{
					Message: "User not found",
				})
				tools.HandleError(err)
			case tools.ForumExist:
				err := c.JSON(http.StatusConflict, f)
				tools.HandleError(err)
			default:
				logrus.Error(err)
				return err
			}
		} else {
			err = c.JSON(http.StatusCreated, f)
			tools.HandleError(err)
		}
		return nil
	}
}

func (fh *ForumHandler) GetForumInfo() echo.HandlerFunc {
	return func(c echo.Context) error {
		/*logrus.WithFields(logrus.Fields{
			"method": c.Request().Method,
		}).Info(c.Request().URL)*/
		f := &models.Forum{}
		f.Slug = c.Param("slug")
		if err := fh.uc.GetForum(f); err != nil {
			err := c.JSON(http.StatusNotFound, tools.Message{
				Message: "forum not found",
			})
			tools.HandleError(err)
			return nil
		}
		err := c.JSON(http.StatusOK, f)
		tools.HandleError(err)
		return nil
	}
}

func (fh *ForumHandler) GetForumThreads() echo.HandlerFunc {
	return func(c echo.Context) error {
		/*logrus.WithFields(logrus.Fields{
			"method": c.Request().Method,
		}).Info(c.Request().URL)*/
		f := &models.Forum{}
		f.Slug = c.Param("slug")
		err := c.Bind(f)
		tools.HandleError(err)
		ths, err := fh.uc.GetForumThreads(
			f,
			c.QueryParam("desc"),
			c.QueryParam("limit"),
			c.QueryParam("since"),
		)
		if err != nil {
			switch err {
			case tools.ForumNotExist:
				err := c.JSON(http.StatusNotFound,
					tools.Message{
						Message: "forum not found",
					})
				tools.HandleError(err)
			default:
				tools.HandleError(err)
				return err
			}
		} else {
			tools.HandleError(err)
			err = c.JSON(http.StatusOK, ths)
		}
		return nil
	}
}

func (fh *ForumHandler) GetForumUsers() echo.HandlerFunc {
	return func(c echo.Context) error {
		/*logrus.WithFields(logrus.Fields{
			"method": c.Request().Method,
		}).Info(c.Request().URL)*/
		f := &models.Forum{}
		f.Slug = c.Param("slug")
		usrs, err := fh.uc.GetForumUsers(f, c.QueryParam("desc"), c.QueryParam("limit"), c.QueryParam("since"))
		if err != nil {
			switch err {
			case tools.ForumNotExist:
				err := c.JSON(http.StatusNotFound,
					tools.Message{
						Message: "forum not found",
					})
				tools.HandleError(err)
			default:
				tools.HandleError(err)
				return err
			}
		} else {
			err = c.JSON(http.StatusOK, usrs)
			tools.HandleError(err)
		}
		return nil
	}
}
