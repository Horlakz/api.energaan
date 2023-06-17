package app

import (
	"net/http"

	"github.com/gofiber/fiber/v2"

	dto "github.com/horlakz/energaan-api/database/dto/app"
	services "github.com/horlakz/energaan-api/database/services/app"
	userService "github.com/horlakz/energaan-api/database/services/auth"
	"github.com/horlakz/energaan-api/handler"
	"github.com/horlakz/energaan-api/helper"
	categoryRequest "github.com/horlakz/energaan-api/payload/request/app/category"
	"github.com/horlakz/energaan-api/payload/response"
	validators "github.com/horlakz/energaan-api/validator/app"
)

type CategoryResult struct {
	Category         dto.CategoryDTO `json:"category"`
	CreatedByDetails string          `json:"createdByDetails"`
}

type CategoryHandlerInterface interface {
	IndexHandle(c *fiber.Ctx) error
	CreateHandle(c *fiber.Ctx) error
	UpdateHandle(c *fiber.Ctx) error
	DeleteHandle(c *fiber.Ctx) error
}

type categoryHandler struct {
	handler.BaseHandler
	categoryService   services.CategoryServiceInterface
	userService       userService.UserServiceInterface
	categoryValidator validators.CategoryValidator
}

func NewCategoryHandler(
	categoryService services.CategoryServiceInterface,
	userService userService.UserServiceInterface,
) CategoryHandlerInterface {
	return &categoryHandler{
		categoryService: categoryService,
		userService:     userService,
	}
}

func (handler *categoryHandler) IndexHandle(c *fiber.Ctx) (err error) {
	var resp response.Response
	categories, queryError := handler.categoryService.ReadAll()

	if queryError != nil {
		resp.Status = http.StatusUnprocessableEntity
		resp.Message = queryError.Error()

		return c.Status(http.StatusUnprocessableEntity).JSON(resp)
	}

	results := make([]CategoryResult, len(categories))

	for i, category := range categories {
		createdByDetails, err := handler.userService.Read(category.CreatedByID)

		if err != nil {
			resp.Status = http.StatusUnprocessableEntity
			resp.Message = err.Error()

			return c.Status(http.StatusUnprocessableEntity).JSON(resp)
		}

		results[i] = CategoryResult{
			Category:         category,
			CreatedByDetails: createdByDetails.FullName,
		}
	}

	resp.Status = http.StatusOK
	resp.Message = http.StatusText(http.StatusOK)

	if len(categories) < 1 {
		resp.Data = map[string]interface{}{"result": []string{}}
		return c.Status(http.StatusOK).JSON(resp)
	}

	resp.Data = map[string]interface{}{"result": results}
	return c.Status(http.StatusOK).JSON(resp)

}

func (handler *categoryHandler) CreateHandle(c *fiber.Ctx) (err error) {
	var resp response.Response
	var createRequest categoryRequest.CreateRequest

	userID, _ := handler.GetUserID(c)
	categoryDTO := new(dto.CategoryDTO)

	// parse request body to dto
	if err := c.BodyParser(&createRequest); err != nil {
		resp.Status = http.StatusExpectationFailed
		resp.Message = err.Error()

		return c.Status(fiber.StatusBadRequest).JSON(resp)
	}

	categoryDTO.Name = createRequest.Name
	categoryDTO.Slug = helper.CreateSlug(categoryDTO.Name)
	categoryDTO.CreatedByID = userID

	_, existErr := handler.categoryService.Read(categoryDTO.Slug)

	if existErr == nil {
		resp.Status = http.StatusConflict
		resp.Message = "Category Already Exist"

		return c.Status(http.StatusConflict).JSON(resp)
	}

	vEs, err := handler.categoryValidator.Validate(*categoryDTO)

	if err != nil {
		resp.Status = http.StatusUnprocessableEntity
		resp.Message = "validation Error"
		resp.Data = vEs

		return c.Status(http.StatusUnprocessableEntity).JSON(resp)
	}

	if _, err := handler.categoryService.Create(*categoryDTO); err != nil {
		resp.Status = http.StatusInternalServerError
		resp.Message = err.Error()

		return c.Status(http.StatusInternalServerError).JSON(resp)
	}

	resp.Status = http.StatusCreated
	resp.Message = http.StatusText(http.StatusCreated)
	resp.Data = map[string]interface{}{"result": categoryDTO}

	return c.Status(http.StatusOK).JSON(resp)
}

func (handler *categoryHandler) UpdateHandle(c *fiber.Ctx) (err error) {
	var resp response.Response
	userID, err := handler.GetUserID(c)

	slug := c.Params("slug")

	if err != nil {
		resp.Status = http.StatusExpectationFailed
		resp.Message = "Exception Error: " + err.Error()

		return c.Status(http.StatusExpectationFailed).JSON(resp)
	}

	categoryDto, err := handler.categoryService.Read(slug)
	categoryDto.Slug = helper.CreateSlug(categoryDto.Name)

	if err != nil {
		resp.Status = http.StatusNotFound
		resp.Message = "Record Not Found"

		return c.Status(http.StatusNotFound).JSON(resp)
	}

	if err := c.BodyParser(&categoryDto); err != nil {
		resp.Status = http.StatusExpectationFailed
		resp.Message = err.Error()

		return c.Status(http.StatusExpectationFailed).JSON(resp)
	}

	categoryDto.Slug = slug
	categoryDto.UpdatedByID = userID

	vEs, err := handler.categoryValidator.Validate(categoryDto)

	if err != nil {
		resp.Status = http.StatusUnprocessableEntity
		resp.Message = "validation Error"
		resp.Data = vEs

		return c.Status(http.StatusUnprocessableEntity).JSON(resp)
	}

	if _, err := handler.categoryService.Update(categoryDto); err != nil {
		resp.Status = http.StatusInternalServerError
		resp.Message = err.Error()

		return c.Status(http.StatusInternalServerError).JSON(resp)
	}

	resp.Status = 200
	resp.Message = "Update Successful"

	return c.Status(200).JSON(resp)
}

func (handler *categoryHandler) DeleteHandle(c *fiber.Ctx) (err error) {

	var resp response.Response
	userID, err := handler.GetUserID(c)

	slug := c.Params("slug")

	if err != nil {
		resp.Status = http.StatusExpectationFailed
		resp.Message = "Exception Error: " + err.Error()

		return c.Status(http.StatusExpectationFailed).JSON(resp)
	}

	_, existsErr := handler.categoryService.Read(slug)

	if existsErr != nil {
		resp.Status = http.StatusNotFound
		resp.Message = "Record not found"
		return c.Status(http.StatusNotFound).JSON(resp)
	}

	if err := handler.categoryService.Delete(userID, slug); err != nil {
		resp.Status = http.StatusInternalServerError
		resp.Message = err.Error()

		return c.Status(http.StatusInternalServerError).JSON(resp)
	}

	resp.Status = http.StatusOK
	resp.Message = http.StatusText(http.StatusOK)

	return c.Status(http.StatusOK).JSON(resp)
}
