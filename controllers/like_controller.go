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
	var existingLike models.Like
	err = database.DB.Where("user_id = ? AND blog_id = ?", UserID, BlogID).Take(&existingLike).Error
	if err == nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Already liked the blog"})
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

func UnlikeBlog(c *fiber.Ctx) error {
	UserID, ok := c.Locals("user").(jwt.MapClaims)["sub"]
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Token not valid"})
	}

	BlogID, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid 'id' parameter"})
	}

	// Check if the like exists
	var existingLike models.Like
	err = database.DB.Where("user_id = ? AND blog_id = ?", UserID, BlogID).Take(&existingLike).Error
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "User has not liked the blog"})
	}

	// Delete the like
	err = database.DB.Delete(&existingLike).Error
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Blog unliked successfully"})
}
