package handler

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	"github.com/horlakz/energaan-api/database/repository"
	"github.com/horlakz/energaan-api/helper"
	"github.com/horlakz/energaan-api/payload/response"
	"github.com/horlakz/energaan-api/security"
)

type BaseHandler struct{}

func Index(c *fiber.Ctx) error {

	var resp response.Response

	var about struct {
		Name    string `json:"name"`
		Version string `json:"version"`
		Author  string `json:"author"`
	}

	about.Name = "Energaan API"
	about.Version = "0.0.2"
	about.Author = "Horlakz"

	resp.Status = http.StatusOK
	resp.Message = http.StatusText(http.StatusOK)
	resp.Data = map[string]interface{}{"about": about}

	return c.JSON(resp)
}

func NotFound(c *fiber.Ctx) error {
	var resp response.Response

	resp.Status = http.StatusNotFound
	resp.Message = "Route not found"

	return c.Status(http.StatusNotFound).JSON(resp)
}

func (baseHandler *BaseHandler) GetUserID(c *fiber.Ctx) (uuid.UUID, error) {
	userID, err := security.ExtractUserID(c.Request())

	if err != nil {
		return uuid.Nil, err
	}

	return userID, nil
}

func (h *BaseHandler) GeneratePageable(context *fiber.Ctx) (pageable repository.Pageable) {

	pageable.Page = 1
	pageable.Size = 20
	pageable.SortBy = "created_at"
	pageable.SortDirection = "asc"
	pageable.Search = ""

	size, err := strconv.Atoi(context.Query("size", "0"))
	if (size > 0) && err == nil {
		pageable.Size = size
	}

	page, err := strconv.Atoi(context.Query("page", "1"))
	if (page > 0) && err == nil {
		pageable.Page = page
	}

	orderBy := context.Query("sort_by", "")
	if orderBy != "" {
		pageable.SortBy = orderBy
	}

	sortDir := context.Query("sort_dir", "")
	if sortDir != "" {
		pageable.SortBy = sortDir
	}

	search := context.Query("search", "")
	if search != "" {
		pageable.Search = search
	}

	return pageable
}

func StreamFile(c *fiber.Ctx) error {
	file := c.Params("file")
	return c.SendFile("./images/" + file)
}

func StreamFileFromAwsS3(c *fiber.Ctx) error {
	obj := c.Params("key")
	var resp response.Response

	media := helper.NewMediaHelper()

	file, getErr := media.GetObjectFromS3(obj)

	if getErr != nil {
		log.Println("Error getting S3 file:", getErr)
		resp.Status = http.StatusUnprocessableEntity
		resp.Message = fmt.Sprintf("Error getting S3 file: %s", getErr.Error())
		return c.Status(http.StatusNotFound).JSON(resp)
	}

	c.Set(fiber.HeaderContentType, *file.ContentType)
	c.Set(fiber.HeaderContentLength, strconv.FormatInt(*file.ContentLength, 10))

	_, err := io.Copy(c, file.Body)

	if err != nil {
		resp.Status = http.StatusUnprocessableEntity
		resp.Message = fmt.Sprintf("Error streaming S3 file: %s", err.Error())
		return c.Status(http.StatusNotFound).JSON(resp)
	}

	return nil
}
