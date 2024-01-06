package main

import (
	"golang_project_1/database"
	"golang_project_1/routers"

	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/joho/godotenv"
)

func init() {
	//env
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	//dabase
	database.InitDB()
}
func main() {
	app := fiber.New()
	
	// middlewares
	app.Use(recover.New())
	app.Use(cors.New())
	app.Use(logger.New())
	app.Static("/", "./public")

	// routers
	routers.SetupRouters(app)

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})
	
	log.Fatal(app.Listen(os.Getenv("PORT")))

}
