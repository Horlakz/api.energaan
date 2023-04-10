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

func (m *Media) Save(c *fiber.Ctx) (string, error) {
	form, formErr := c.MultipartForm()

	if formErr != nil {
		return "", formErr
	}

	// get image file from form
	file := form.File["image"][0]

	fileName := m.FileName(file.Filename)

	if err := c.SaveFile(file, fmt.Sprintf("./images/%s", fileName)); err != nil {
		return "", err
	}

	return fileName, nil
}
