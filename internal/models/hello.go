package models

type Hello struct {
	ID      int64  `json:"id"`
	Content string `json:"content"`
}

func (Hello) TableName() string {
	return "hello"
}
