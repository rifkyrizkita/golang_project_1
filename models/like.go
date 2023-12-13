package models

type Like struct {
	ID     uint
	BlogID uint
	Blog   Blog
	UserID uint
	User   User
}
