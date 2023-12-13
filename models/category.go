package models

type Category struct {
	ID       uint
	Category string `gorm:"type:varchar(255);not null;unique"`
	Blog     *Blog  `json:"-"`
}
