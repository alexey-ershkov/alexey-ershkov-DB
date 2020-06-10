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
		//logrus.WithFields(logrus.Fields{
		//	"method": c.Request().Method,
		//}).Info(c.Request().URL)
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

///api/forum/-O6MJSPR6XCMR/threads?desc=true&limit=16&since=2020-11-11T20%3A13%3A17.178%2B03%3A00  method=GET
//INFO[0043] /api/forum/iA-Mf848OvjMr/threads?desc=true&limit=19&since=2020-10-07T09%3A22%3A29.561%2B03%3A00  method=GET
//INFO[0043] /api/forum/0V-MJSP86XJMS/threads?limit=16     method=GET
//INFO[0043] /api/forum/pmoLC84R6E53s/threads?desc=true&limit=15&since=2019-12-29T01%3A55%3A42.339%2B03%3A00  method=GET
//INFO[0043] /api/forum/sj63580k6X5Ls/threads?desc=true&limit=15&since=2021-04-18T10%3A06%3A03.148%2B03%3A00  method=GET
