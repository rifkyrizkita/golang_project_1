package routers

import (
	"golang_project_1/controllers"
	"golang_project_1/middlewares"

	"github.com/gofiber/fiber/v2"
)

func LikeRouters(like fiber.Router) {
	// post routers
	like.Post("/:id", middlewares.VerifyToken, controllers.LikeBlog)
	// delete routers
	like.Delete("/:id",middlewares.VerifyToken ,controllers.UnlikeBlog)
}
