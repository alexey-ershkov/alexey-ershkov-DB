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
)
