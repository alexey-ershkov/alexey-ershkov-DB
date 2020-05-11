package models

import "github.com/jackc/pgtype"

type Post struct {
	Author   string           `json:"author"`
	Created  string           `json:"created"`
	Forum    string           `json:"forum",url:"param"`
	Id       int64            `json:"id"`
	IsEdited bool             `json:"isEdited"`
	Message  string           `json:"message"`
	Parent   int64            `json:"parent"`
	Thread   int64            `json:"thread"`
	Path     pgtype.Int8Array `json:"-"`
}
