package controllers

import (
	"golang_project_1/database"
	"golang_project_1/helpers"
	"golang_project_1/models"
	"golang_project_1/web/requests"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

func CreateBlog(c *fiber.Ctx) error {
	id, ok := c.Locals("user").(jwt.MapClaims)["sub"].(float64)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Token not valid"})
	}

	filename := c.Locals("filename").(string)

	var body requests.BlogCreateBody
	err := c.BodyParser(&body)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	blog := models.Blog{
		Title:      body.Title,
		Content:    body.Content,
		ImageBlog:  filename,
		CategoryID: body.CategoryID,
		UserID:     uint(id),
	}
	err = database.DB.Create(&blog).Error
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	blog = models.Blog{}
	err = database.DB.Joins("User").Joins("Category").Last(&blog).Error
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": err.Error()})
	}
	response := fiber.Map{
		"message": "Blog post created successfully!",
		"result":  blog,
	}

	return c.Status(fiber.StatusCreated).JSON(response)
}

func FindById(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid 'id' parameter"})
	}
	var blog models.Blog

	err = database.DB.Joins("User").Joins("Category").Where("blogs.id = ?", id).Take(&blog).Error
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": err.Error()})
	}
	var likesCount struct {
		Count int
	}

	if err := database.DB.Model(&models.Like{}).
		Select("COUNT(*) as count").
		Where("blog_id = ?", id).
		Scan(&likesCount).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	// Menggabungkan total likes ke dalam hasil blog
	blog.TotalLikes = likesCount.Count

	response := fiber.Map{
		"result": blog,
	}

	return c.Status(fiber.StatusOK).JSON(response)
}
func FindAll(c *fiber.Ctx) error {
	page := c.QueryInt("page", 1)
	limit := c.QueryInt("limit", 5)
	search := c.Query("search")
	order := c.Query("order", "desc")
	offset := (page - 1) * limit

	var blogs []models.Blog
	if search != "" {
		err := database.DB.Joins("User").
			Joins("Category").
			Where("title like ?", "%"+search+"%").
			Limit(limit).
			Offset(offset).
			Order("created_at " + order).
			Find(&blogs).
			Error

		if err != nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": err.Error()})
		}
	} else {
		err := database.DB.Joins("User").
			Joins("Category").
			Limit(limit).
			Offset(offset).
			Order("created_at " + order).
			Find(&blogs).
			Error

		if err != nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": err.Error()})
		}
	}
	var likesCount []struct {
		BlogID uint
		Count  int
	}
	if err := database.DB.Model(&models.Like{}).
		Select("blog_id, COUNT(*) as count").
		Group("blog_id").
		Scan(&likesCount).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	for i, blog := range blogs {
		for _, likeCount := range likesCount {
			if blog.ID == likeCount.BlogID {
				blogs[i].TotalLikes = likeCount.Count
				break
			}
		}
	}

	response := fiber.Map{
		"result": blogs,
	}

	return c.Status(fiber.StatusOK).JSON(response)
}

func FindByUserID(c *fiber.Ctx) error {
	page := c.QueryInt("page", 1)
	limit := c.QueryInt("limit", 5)
	search := c.Query("search")
	order := c.Query("order", "desc")
	offset := (page - 1) * limit

	id, ok := c.Locals("user").(jwt.MapClaims)["sub"]
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Token not valid"})
	}

	var blogs []models.Blog

	if search != "" {
		err := database.DB.Joins("User").
			Joins("Category").
			Where("title like ?", "%"+search+"%").
			Where("user_id = ?", id).
			Limit(limit).
			Offset(offset).
			Order("created_at " + order).
			Find(&blogs).
			Error

		if err != nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": err.Error()})
		}
	} else {
		err := database.DB.Joins("User").
			Joins("Category").
			Where("user_id = ?", id).
			Limit(limit).
			Offset(offset).
			Order("created_at " + order).
			Find(&blogs).
			Error

		if err != nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": err.Error()})
		}
	}
	var likesCount []struct {
		BlogID uint
		Count  int
	}
	if err := database.DB.Model(&models.Like{}).
		Select("blog_id, COUNT(*) as count").
		Where("blog_id IN ?", helpers.GetBlogIDs(&blogs)).
		Group("blog_id").
		Scan(&likesCount).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	// Menggabungkan total likes ke dalam hasil blog
	for i, blog := range blogs {
		for _, likeCount := range likesCount {
			if blog.ID == likeCount.BlogID {
				blogs[i].TotalLikes = likeCount.Count
				break
			}
		}
	}

	response := fiber.Map{
		"result": blogs,
	}

	return c.Status(fiber.StatusOK).JSON(response)
}

func DeleteBlog(c *fiber.Ctx) error {
	UserID, ok := c.Locals("user").(jwt.MapClaims)["sub"].(float64)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "User not authenticated"})
	}

	BlogID, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid 'id' parameter"})
	}

	var blog models.Blog
	err = database.DB.Where("id = ?", BlogID).Take(&blog).Error
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Blog not found"})
	}

	if uint(UserID) != blog.UserID {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized: Wrong user id"})
	}

	err = database.DB.Delete(&blog, BlogID).Error
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	response := fiber.Map{
		"message": "Blog deleted",
	}

	return c.Status(fiber.StatusOK).JSON(response)
}
