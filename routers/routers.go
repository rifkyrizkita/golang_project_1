// routes.go
package routers

import (
	"github.com/gofiber/fiber/v2"
)

func SetupRouters(app *fiber.App) {
	api := app.Group("/api")

	user := api.Group("/user")
	UserRouters(user)

	blog := api.Group("/blog")
	BlogRouters(blog)

	like := api.Group("/like")
	LikeRouters(like)

	test := api.Group("/test")
	TesterRouters(test)
}
