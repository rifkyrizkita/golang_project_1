package helpers

import "golang_project_1/models"

func GetBlogIDs(blogs *[]models.Blog) []interface{} {
	var blogIDs []interface{}
	for _, blog := range *blogs {
		blogIDs = append(blogIDs, blog.ID)
	}
	return blogIDs
}
