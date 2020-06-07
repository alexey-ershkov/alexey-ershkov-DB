package delivery

import (
	"alexey-ershkov/alexey-ershkov-DB.git/internal/forum"
	"alexey-ershkov/alexey-ershkov-DB.git/internal/models"
	"alexey-ershkov/alexey-ershkov-DB.git/internal/post"
	"alexey-ershkov/alexey-ershkov-DB.git/internal/thread"
	"alexey-ershkov/alexey-ershkov-DB.git/internal/tools"
	"alexey-ershkov/alexey-ershkov-DB.git/internal/user"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"net/http"
	"strconv"
	"strings"
)

type PostHandler struct {
	pUC  post.Usecase
	fUC  forum.Usecase
	uUC  user.Usecase
	thUC thread.Usecase
}

func Find(slice []string, val string) bool {
	for _, item := range slice {
		if item == val {
			return true
		}
	}
	return false
}

func NewPostHandler(router *echo.Echo, pUC post.Usecase, fUC forum.Usecase, uUC user.Usecase, thUC thread.Usecase) {
	ph := &PostHandler{
		pUC:  pUC,
		uUC:  uUC,
		fUC:  fUC,
		thUC: thUC,
	}

	router.POST("/api/thread/:slug/create", ph.CreatePosts())
	router.GET("/api/post/:id/details", ph.GetPost())
	router.POST("/api/post/:id/details", ph.UpdatePost())
}

func (ph *PostHandler) CreatePosts() echo.HandlerFunc {
	return func(c echo.Context) error {
		/*logrus.WithFields(logrus.Fields{
			"method": c.Request().Method,
		}).Info(c.Request().URL)*/
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
			fmt.Println(e)
			fmt.Println(err)
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
		/*logrus.WithFields(logrus.Fields{
			"method": c.Request().Method,
		}).Info(c.Request().URL)*/
		str := c.QueryParam("related")
		related := strings.Split(str, ",")
		p := &models.Post{}
		th := &models.Thread{}
		f := &models.Forum{}
		u := &models.User{}
		var err error
		p.Id, err = strconv.ParseInt(c.Param("id"), 10, 64)
		if err != nil {
			logrus.Error("Cannot parse id")
		}
		if err := ph.pUC.GetPost(p); err != nil {
			err := c.JSON(http.StatusNotFound, tools.Message{
				Message: err.Error(),
			})
			tools.HandleError(err)
			return nil
		}
		if Find(related, "user") {
			u.Nickname = p.Author
			err := ph.uUC.GetUser(u)
			tools.HandleError(err)
		} else {
			u = nil
		}
		if Find(related, "thread") {
			th.Slug = strconv.FormatInt(p.Thread, 10)
			err := ph.thUC.GetThreadInfo(th)
			tools.HandleError(err)
		} else {
			th = nil
		}
		if Find(related, "forum") {
			f.Slug = p.Forum
			err := ph.fUC.GetForum(f)
			tools.HandleError(err)
		} else {
			f = nil
		}
		ans := models.Info{
			Post:   p,
			Forum:  f,
			Thread: th,
			Author: u,
		}
		err = c.JSON(http.StatusOK, ans)
		tools.HandleError(err)
		return nil
	}
}

func (ph *PostHandler) UpdatePost() echo.HandlerFunc {
	return func(c echo.Context) error {
		/*logrus.WithFields(logrus.Fields{
			"method": c.Request().Method,
		}).Info(c.Request().URL)*/
		p := &models.Post{}
		err := c.Bind(p)
		tools.HandleError(err)
		p.Id, err = strconv.ParseInt(c.Param("id"), 10, 64)
		tools.HandleError(err)
		if err := ph.pUC.UpdatePost(p); err != nil {
			err = c.JSON(http.StatusNotFound, tools.Message{
				Message: "post not found",
			})
			return nil
		}
		err = c.JSON(http.StatusOK, p)
		tools.HandleError(err)
		return nil
	}
}
