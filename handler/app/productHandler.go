package app

import (
	"fmt"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	dto "github.com/horlakz/energaan-api/database/dto/app"
	services "github.com/horlakz/energaan-api/database/services/app"
	"github.com/horlakz/energaan-api/handler"
	"github.com/horlakz/energaan-api/helper"
	"github.com/horlakz/energaan-api/payload/response"
	validators "github.com/horlakz/energaan-api/validator/app"
)

type ProductHandlerInterface interface {
	IndexHandle(c *fiber.Ctx) error
	CreateHandle(c *fiber.Ctx) error
	ReadHandle(c *fiber.Ctx) error
	UpdateHandle(c *fiber.Ctx) error
	DeleteHandle(c *fiber.Ctx) error
}

type productHandler struct {
	handler.BaseHandler
	mediaHelper      helper.MediaInterface
	productService   services.ProductServiceInterface
	productValidator validators.ProductValidator
}

func NewProductHandler(productService services.ProductServiceInterface) ProductHandlerInterface {
	return &productHandler{
		productService: productService,
	}
}

func (handler *productHandler) IndexHandle(c *fiber.Ctx) (err error) {
	var resp response.Response
	pageable := handler.GeneratePageable(c)
	categoryId := c.Query("categoryId")

	var category uuid.UUID

	if categoryId != "" {
		category, err = uuid.Parse(categoryId)
		if err != nil {
			resp.Status = http.StatusUnprocessableEntity
			resp.Message = err.Error()
		}
	}

	products, pagination, queryError := handler.productService.ReadAll(pageable, category)

	if queryError != nil {
		resp.Status = http.StatusUnprocessableEntity
		resp.Message = queryError.Error()

		return c.Status(http.StatusUnprocessableEntity).JSON(resp)
	}

	resp.Status = http.StatusOK
	resp.Message = http.StatusText(http.StatusOK)
	resp.Data = map[string]interface{}{
		"result":      products,
		"totalPages":  pagination.TotalPages,
		"totalItems":  pagination.TotalItems,
		"currentPage": pagination.CurrentPage,
	}

	return c.Status(http.StatusOK).JSON(resp)
}

func (handler *productHandler) CreateHandle(c *fiber.Ctx) (err error) {
	var resp response.Response

	userID, _ := handler.GetUserID(c)
	productDto := new(dto.ProductDTO)

	// collect fields and image file from request as multipart form
	form, formErr := c.MultipartForm()

	if formErr != nil {
		resp.Status = http.StatusExpectationFailed
		resp.Message = formErr.Error()

		return c.Status(http.StatusBadRequest).JSON(resp)
	}

	// parse request body to dto
	if err := c.BodyParser(productDto); err != nil {
		resp.Status = http.StatusExpectationFailed
		resp.Message = err.Error()

		fmt.Print(err.Error(), "body parser")

		return c.Status(http.StatusBadRequest).JSON(resp)
	}

	fileNames, _ := handler.mediaHelper.Save(c)

	for _, filename := range fileNames {
		err := handler.mediaHelper.UploadToAWSS3(filename)

		if err != nil {
			resp.Status = http.StatusExpectationFailed
			resp.Message = err.Error()

			return c.Status(http.StatusBadRequest).JSON(resp)
		}
	}

	// save image file to dto
	productDto.Images = fileNames

	// get other fields from form
	productDto.Title = form.Value["title"][0]
	productDto.CategoryID, err = uuid.Parse(form.Value["categoryId"][0])
	productDto.Description = form.Value["description"][0]
	productDto.Features = form.Value["features"]

	productDto.CreatedByID = userID
	productDto.Slug = helper.CreateSlug(productDto.Title)

	if err != nil {
		resp.Status = http.StatusUnprocessableEntity
		resp.Message = err.Error()
	}

	_, existErr := handler.productService.Read(productDto.Slug)

	if existErr == nil {
		resp.Status = http.StatusConflict
		resp.Message = "Record Already Exist"

		return c.Status(http.StatusConflict).JSON(resp)
	}

	vEs, err := handler.productValidator.Validate(*productDto)

	if err != nil {
		resp.Status = http.StatusUnprocessableEntity
		resp.Message = "validation Error"
		resp.Data = vEs

		return c.Status(http.StatusUnprocessableEntity).JSON(resp)
	}

	if _, err := handler.productService.Create(*productDto); err != nil {
		resp.Status = http.StatusInternalServerError
		resp.Message = err.Error()

		return c.Status(http.StatusInternalServerError).JSON(resp)
	}

	resp.Status = http.StatusCreated
	resp.Message = http.StatusText(http.StatusCreated)
	resp.Data = map[string]interface{}{"product": productDto}

	return c.Status(http.StatusOK).JSON(resp)
}

func (handler *productHandler) ReadHandle(c *fiber.Ctx) (err error) {
	var resp response.Response

	slug := c.Params("slug")

	// if err != nil {
	// 	resp.Status = http.StatusExpectationFailed
	// 	resp.Message = "Exception Error: " + err.Error()

	// 	return c.Status(http.StatusExpectationFailed).JSON(resp)
	// }

	productDto, err := handler.productService.Read(slug)

	if err != nil {
		resp.Status = http.StatusNotFound
		resp.Message = "Record Not Found"

		return c.Status(http.StatusNotFound).JSON(resp)
	}

	resp.Status = http.StatusOK
	resp.Message = http.StatusText(http.StatusOK)
	resp.Data = map[string]interface{}{"product": productDto}

	return c.Status(200).JSON(resp)
}

func (handler *productHandler) UpdateHandle(c *fiber.Ctx) (err error) {
	var resp response.Response
	userID, err := handler.GetUserID(c)

	slug := c.Params("slug")

	if err != nil {
		resp.Status = http.StatusExpectationFailed
		resp.Message = "Exception Error: " + err.Error()

		return c.Status(http.StatusExpectationFailed).JSON(resp)
	}

	productDto, err := handler.productService.Read(slug)
	productDto.Slug = helper.CreateSlug(productDto.Title)

	if err != nil {
		resp.Status = http.StatusNotFound
		resp.Message = "Record Not Found"

		return c.Status(http.StatusNotFound).JSON(resp)
	}

	if err := c.BodyParser(&productDto); err != nil {
		resp.Status = http.StatusExpectationFailed
		resp.Message = err.Error()

		return c.Status(http.StatusExpectationFailed).JSON(resp)
	}

	productDto.Slug = slug
	productDto.UpdatedByID = userID

	vEs, err := handler.productValidator.Validate(productDto)

	if err != nil {
		resp.Status = http.StatusUnprocessableEntity
		resp.Message = "validation Error"
		resp.Data = vEs

		return c.Status(http.StatusUnprocessableEntity).JSON(resp)
	}

	if _, err := handler.productService.Update(productDto); err != nil {
		resp.Status = http.StatusInternalServerError
		resp.Message = err.Error()

		return c.Status(http.StatusInternalServerError).JSON(resp)
	}

	resp.Status = 200
	resp.Message = "Update Successful"

	return c.Status(200).JSON(resp)
}

func (handler *productHandler) DeleteHandle(c *fiber.Ctx) (err error) {

	var resp response.Response
	userID, err := handler.GetUserID(c)

	slug := c.Params("slug")

	if err != nil {
		resp.Status = http.StatusExpectationFailed
		resp.Message = "Exception Error: " + err.Error()

		return c.Status(http.StatusExpectationFailed).JSON(resp)
	}

	_, existsErr := handler.productService.Read(slug)

	if existsErr != nil {
		resp.Status = http.StatusNotFound
		resp.Message = "Record not found"
		return c.Status(http.StatusNotFound).JSON(resp)
	}

	if err := handler.productService.Delete(userID, slug); err != nil {
		resp.Status = http.StatusInternalServerError
		resp.Message = err.Error()

		return c.Status(http.StatusInternalServerError).JSON(resp)
	}

	resp.Status = http.StatusOK
	resp.Message = http.StatusText(http.StatusOK)

	return c.Status(http.StatusOK).JSON(resp)
}
