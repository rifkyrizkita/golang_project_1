package models

import "gorm.io/gorm"

type Blog struct {
	gorm.Model
	Title      string `gorm:"type:varchar(255);not null;unique"`
	Content    string `gorm:"not null"`
	ImageBlog  string `gorm:"type:varchar(255)"`
	TotalLikes int    `gorm:"-"`
	CategoryID uint
	Category   Category
	UserID     uint
	User       User
}
