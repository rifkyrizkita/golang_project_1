package main

import (
	"golang_project_1/database"
	"golang_project_1/models"
	"log"

	"github.com/joho/godotenv"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	database.InitDB()
}

func main() {
	database.DB.AutoMigrate(&models.User{}, &models.Blog{}, &models.Category{})
	database.DB.AutoMigrate(&models.Like{})
}
