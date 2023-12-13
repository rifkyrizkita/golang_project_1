package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Username     string  `gorm:"type:varchar(255);not null;unique"`
	Email        string  `gorm:"type:varchar(255);not null;unique"`
	Phone        string  `gorm:"type:varchar(255);not null;unique"`
	Password     string  `gorm:"type:varchar(255);not null"`
	ImageProfile *string `gorm:"type:varchar(255)"`
	IsVerified   bool    `gorm:"type:boolean;default:false"`
	Blogs        []Blog  `json:"-"`
}
