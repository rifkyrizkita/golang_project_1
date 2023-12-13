package helpers

import (
	"bytes"
	"html/template"
	"os"

	"github.com/gofiber/fiber/v2"
	"gopkg.in/gomail.v2"
)

func SendEmailWithHTMLFile(c *fiber.Ctx, to, subject, htmlFilePath string, data interface{}) error {

	var bodyBuffer bytes.Buffer
	t, err := template.ParseFiles(htmlFilePath)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	err = t.Execute(&bodyBuffer, data)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	m := gomail.NewMessage()
	m.SetHeader("From", os.Getenv("MY_EMAIL_ADDRESS"))
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", bodyBuffer.String())

	d := gomail.NewDialer("smtp.gmail.com", 587, os.Getenv("MY_EMAIL_ADDRESS"), os.Getenv("MY_EMAIL_PASS"))

	if err := d.DialAndSend(m); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return nil
}
