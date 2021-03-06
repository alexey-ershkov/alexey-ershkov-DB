package models

type Thread struct {
	Author  string `json:"author"`
	Created string `json:"created,omitempty"`
	Forum   string `json:"forum"`
	Id      int64  `json:"id"`
	Message string `json:"message"`
	Slug    string `json:"slug,omitempty"`
	Title   string `json:"title"`
	Votes   int64  `json:"votes"`
}
