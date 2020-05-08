package tools

import (
	"errors"
	"github.com/sirupsen/logrus"
)

func HandleError(e error) {
	if e != nil {
		logrus.Error(e)
	}
}

var (
	UserExist = errors.New("user exist")
	//ForumExist = errors.New("forum exist")
	ThreadExist    = errors.New("thread already exist")
	UserNotExist   = errors.New("user does't exist")
	ThreadNotExist = errors.New("no such thread")
	SqlError       = errors.New("sql error")
	ParentNotExist = errors.New("post parent not exist")
	PostNotExist   = errors.New("post not exist")
)
