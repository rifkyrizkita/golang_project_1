package routers

import (
	"golang_project_1/controllers"
	"golang_project_1/middlewares"

	"github.com/gofiber/fiber/v2"
)

func LikeRouters(like fiber.Router) {
	//post router
	like.Post("/:id", middlewares.VerifyToken, controllers.LikeBlog)
}
