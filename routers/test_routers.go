package routers

import (
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
)

func TesterRouters(test fiber.Router) {

	// upload multiple files
	test.Post("/upload", func(c *fiber.Ctx) error {
		form, err := c.MultipartForm()
		if err != nil {
			return err
		}

		files := form.File["documents"]
		var successMessages []string

		for _, file := range files {
			filename := strconv.Itoa(rand.Intn(1000)) + filepath.Ext(file.Filename)
			if err := c.SaveFile(file, "public/"+filename); err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
			}
			successMessages = append(successMessages, filename+" uploaded successfully")

		}

		return c.SendString(fmt.Sprintf("All file(s) uploaded successfully: %s", strings.Join(successMessages, ", ")))
	})

	//delete file
	test.Delete("/delete-file/*", func(c *fiber.Ctx) error {
		// Get the filename from the URL parameters
		filename := c.Params("*")

		// Attempt to delete the file
		err := os.Remove(filename)
		if err != nil {
			// If there was an error, return an error response
			return c.Status(fiber.StatusInternalServerError).SendString(fmt.Sprintf("Error deleting file: %s", err.Error()))

		}

		// If successful, return a success response
		file := strings.Split(filename, "/")
		return c.SendString(fmt.Sprintf("File %s deleted successfully", file[len(file)-1]))
	})

}
