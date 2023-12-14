package controllers

import (
	"golang_project_1/database"
	"golang_project_1/helpers"
	"golang_project_1/models"
	"golang_project_1/web/requests"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

func Register(c *fiber.Ctx) error {
	var body requests.RegisterBody
	err := c.BodyParser(&body)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(body.Password), 10)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	user := models.User{
		Username: body.Username,
		Email:    body.Email,
		Phone:    body.Phone,
		Password: string(hashedPassword),
	}

	err = database.DB.Create(&user).Error
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	tokenString, _ := helpers.GenerateToken(user.ID, time.Now().Add(time.Hour*24), c)

	go func() {
		helpers.SendEmailWithHTMLFile(
			c,
			user.Email,
			"Registration Successful",
			"templates/email_verification.html",
			fiber.Map{
				"Username": user.Username,
				"Token":    tokenString,
			},
		)
	}()

	response := fiber.Map{
		"message": "Registration successful. A verification email has been sent to your email address.",
		"result":  user,
		"token":   tokenString,
	}

	return c.Status(fiber.StatusCreated).JSON(response)
}

func Verification(c *fiber.Ctx) error {
	id, ok := c.Locals("user").(jwt.MapClaims)["sub"]
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Token not valid"})
	}

	var user models.User
	database.DB.Where("id = ?", id).Take(&user)
	if user.IsVerified {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Account already verified"})
	}
	user.IsVerified = true
	err := database.DB.Save(&user).Error
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Account successfully verified"})
}

func AllUsers(c *fiber.Ctx) error {
	var users []models.User
	err := database.DB.Find(&users).Error
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"result": users})
}

func Login(c *fiber.Ctx) error {
	var body requests.LoginBody
	err := c.BodyParser(&body)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	var user models.User
	err = database.DB.Where("(username = ? or email = ? or phone = ?)",
		body.Username,
		body.Email,
		body.Phone).Take(&user).Error
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid credentials 1"})
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password))
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid credentials 2"})
	}

	tokenString, _ := helpers.GenerateToken(user.ID, time.Now().Add(time.Hour*24), c)
	c.Cookie(&fiber.Cookie{
		Name:    "token",
		Value:   tokenString,
		Expires: time.Now().Add(1 * time.Hour),
	})

	response := fiber.Map{
		"message": "Login successful",
		"result":  user,
		"token":   tokenString,
	}
	return c.Status(fiber.StatusOK).JSON(response)
}

func Validation(c *fiber.Ctx) error {
	id, ok := c.Locals("user").(jwt.MapClaims)["sub"]
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Token not valid"})
	}

	var user models.User
	err := database.DB.Where("id = ?", id).Take(&user).Error
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "User not found"})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"result": user})
}

func UpdateProfile(c *fiber.Ctx) error {
	id, ok := c.Locals("user").(jwt.MapClaims)["sub"]
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Token not valid"})
	}

	var body requests.UpdateProfileBody
	err := c.BodyParser(&body)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	updateFields := map[string]interface{}{}
	if body.Username != "" {
		updateFields["username"] = body.Username
	}
	if body.Email != "" {
		updateFields["email"] = body.Email
	}
	if body.Phone != "" {
		updateFields["phone"] = body.Phone
	}

	var user models.User
	database.DB.Where("id = ?", id).Take(&user)
	err = database.DB.Model(&user).Updates(updateFields).Error
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Profile updated successfully"})
}
func UpdatePassword(c *fiber.Ctx) error {
	id, ok := c.Locals("user").(jwt.MapClaims)["sub"]
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Token not valid"})
	}

	var body requests.UpdatePasswordBody
	err := c.BodyParser(&body)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	var user models.User
	database.DB.Where("id = ?", id).Take(&user)
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.CurrentPassword))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(body.Password), 10)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	err = database.DB.Model(&user).Update("password", hashedPassword).Error
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Password updated successfully"})
}

func ProfilePicture(c *fiber.Ctx) error {
	filename := c.Locals("filename")
	id, ok := c.Locals("user").(jwt.MapClaims)["sub"]
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Token not valid"})
	}

	var user models.User
	database.DB.Where("id = ?", id).Take(&user)
	err := database.DB.Model(&user).Update("image_profile", filename).Error
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Profile picture updated successfully"})
}

func ForgetPassword(c *fiber.Ctx) error {
	var body requests.ForgetPasswordBody
	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	var user models.User
	if err := database.DB.Where("email = ?", body.Email).Take(&user).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "User not found"})
	}

	tokenString, _ := helpers.GenerateToken(user.ID, time.Now().Add(time.Hour*24), c)

	go func() {
		helpers.SendEmailWithHTMLFile(
			c,
			user.Email,
			"Reset Password",
			"templates/email_reset_password.html",
			fiber.Map{
				"Username": user.Username,
				"Token":    tokenString,
			},
		)
	}()

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Password reset email sent successfully. Please check your email for further instructions.",
		"token":   tokenString,
	})
}

func ResetPassword(c *fiber.Ctx) error {
	id, ok := c.Locals("user").(jwt.MapClaims)["sub"]
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Token not valid"})
	}

	var body requests.ResetPasswordBody
	err := c.BodyParser(&body)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	var user models.User
	database.DB.Where("id = ?", id).Take(&user)

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(body.Password), 10)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	err = database.DB.Model(&user).Update("password", hashedPassword).Error
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Password reset successful. You can now log in with your new password."})
}
