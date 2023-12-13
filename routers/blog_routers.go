package routers

import (
	"golang_project_1/controllers"
	"golang_project_1/middlewares"

	"github.com/gofiber/fiber/v2"
)

func BlogRouters(blog fiber.Router) {
	// post routers
	blog.Post("/", middlewares.VerifyToken, middlewares.UploadFile("BIMG", ""), middlewares.ValidatorCreateBlog, controllers.CreateBlog)
	// get routers
	blog.Get("/", controllers.FindAll)
	blog.Get("/user", middlewares.VerifyToken, controllers.FindByUserID)
	blog.Get("/:id", controllers.FindById)
}
