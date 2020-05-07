package models

type Vote struct {
	Vote     int32  `json:"voice"`
	Nickname string `json:"nickname"`
	Thread   int64  `json:"thread,omitempty"`
}
