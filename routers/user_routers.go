package routers

import (
	"golang_project_1/controllers"
	"golang_project_1/middlewares"

	"github.com/gofiber/fiber/v2"
)

func UserRouters(user fiber.Router) {
	// post routers
	user.Post("/", middlewares.ValidatorRegister, controllers.Register)
	user.Post("/login", middlewares.ValidatorLogin, controllers.Login)
	// patch routers
	user.Patch("/verification-account", middlewares.VerifyToken, controllers.Verification)
	user.Patch("/update-profile", middlewares.VerifyToken, middlewares.ValidatorUpdateProfile, controllers.UpdateProfile)
	user.Patch("/update-password", middlewares.VerifyToken, middlewares.ValidatorUpdatePassword, controllers.UpdatePassword)
	user.Patch("/profile-picture", middlewares.VerifyToken, middlewares.UploadFile("", ""), controllers.ProfilePicture)
	user.Patch("/reset-password", middlewares.VerifyToken, middlewares.ValidatorResetPassword, controllers.ResetPassword)
	// put routers
	user.Put("/forget-password", middlewares.ValidatorForgetPassword, controllers.ForgetPassword)
	// get routers
	user.Get("/all", controllers.AllUsers)
	user.Get("/protected", middlewares.VerifyToken, controllers.Validation)
}
