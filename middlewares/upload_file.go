package middlewares

import (
	"fmt"
	"math/rand"
	"path/filepath"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
)

const maxFileSize = 1 << 20

func UploadFile(prefix string, path string) fiber.Handler {
	if path == "" {
		path = "public/"
	}
	if prefix == "" {
		prefix = "PIMG"
	}

	return func(c *fiber.Ctx) error {
		file, err := c.FormFile("file")
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		}
		
		// Extract the file extension from the original filename
		fileExt := filepath.Ext(file.Filename)

		allowedExts := []string{"jpg", "jpeg", "png", "webp", "gif"}
		ext := strings.ToLower(strings.TrimPrefix(filepath.Ext(file.Filename), "."))

		if !strings.Contains(strings.Join(allowedExts, ","), ext) {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Your file extension is not allowed",
			})
		}

		if file.Size > maxFileSize {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": fmt.Sprintf(
					"File size exceeds the maximum allowed size of 1 megabyte. Actual size: %.2f megabytes",
					float64(file.Size)/float64(1<<20),
				),
			})
		}

		timestamp := time.Now().Unix()

		randomNumber := rand.Intn(100000)

		filename := fmt.Sprintf("%s-%d-%d%s", prefix, timestamp, randomNumber, fileExt)

		err = c.SaveFile(file, filepath.Join(path, filename))
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}

		c.Locals("filename", filename)

		return c.Next()
	}
}
