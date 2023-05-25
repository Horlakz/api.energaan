package helper

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
)

type Media struct{}

func (m *Media) FileName(name string) string {
	postfix := strconv.FormatInt(time.Now().Unix(), 10)
	fileExt := strings.Split(name, ".")[1]
	fileName := strings.Split(name, ".")[0]

	return fileName + postfix + "." + fileExt
}

func (m *Media) Save(c *fiber.Ctx) ([]string, error) {
	form, formErr := c.MultipartForm()

	if formErr != nil {
		return []string{}, formErr
	}

	// get image files from form
	files := form.File["image"]

	if len(files) == 0 {
		return []string{}, fmt.Errorf("image is required")
	}

	fileNames := make([]string, len(files))

	for i, file := range files {
		fileNames[i] = m.FileName(file.Filename)
	}

	for i, file := range files {
		if err := c.SaveFile(file, fmt.Sprintf("./images/%s", fileNames[i])); err != nil {
			return []string{}, err
		}
	}

	return fileNames, nil
}
