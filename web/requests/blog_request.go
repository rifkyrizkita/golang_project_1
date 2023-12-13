package requests

type BlogCreateBody struct {
	Title      string `validate:"required"`
	Content    string `validate:"required"`
	CategoryID uint   `validate:"required,number"`
}