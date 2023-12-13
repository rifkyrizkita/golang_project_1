package controllers

import (
	"golang_project_1/database"
	"golang_project_1/models"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

func LikeBlog(c *fiber.Ctx) error {
	UserID, ok := c.Locals("user").(jwt.MapClaims)["sub"].(float64)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Token not valid"})
	}
	BlogID, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid 'id' parameter"})
	}
	like := models.Like{
		UserID: uint(UserID),
		BlogID: uint(BlogID),
	}
	err = database.DB.Create(&like).Error
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"message": "Blog liked successfully"})
}
