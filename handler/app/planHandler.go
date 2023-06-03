package app

import (
	"net/http"

	"github.com/gofiber/fiber/v2"

	dto "github.com/horlakz/energaan-api/database/dto/app"
	services "github.com/horlakz/energaan-api/database/services/app"
	"github.com/horlakz/energaan-api/handler"
	"github.com/horlakz/energaan-api/helper"
	"github.com/horlakz/energaan-api/payload/response"
	validators "github.com/horlakz/energaan-api/validator/app"
)

type PlanHandlerInterface interface {
	IndexHandle(c *fiber.Ctx) error
	CreateHandle(c *fiber.Ctx) error
	ReadHandle(c *fiber.Ctx) error
	UpdateHandle(c *fiber.Ctx) error
	DeleteHandle(c *fiber.Ctx) error
}

type planHandler struct {
	handler.BaseHandler
	mediaHelper   helper.Media
	planService   services.PlanServiceInterface
	planValidator validators.PlanValidator
}

func NewPlanHandler(planService services.PlanServiceInterface) PlanHandlerInterface {
	return &planHandler{
		planService: planService,
	}
}

func (handler *planHandler) IndexHandle(c *fiber.Ctx) (err error) {
	var resp response.Response
	pageable := handler.GeneratePageable(c)
	plans, pagination, queryError := handler.planService.ReadAll(pageable)

	if queryError != nil {
		resp.Status = http.StatusUnprocessableEntity
		resp.Message = queryError.Error()

		return c.Status(http.StatusUnprocessableEntity).JSON(resp)
	}

	resp.Status = http.StatusOK
	resp.Message = http.StatusText(http.StatusOK)
	resp.Data = map[string]interface{}{"plans": plans, "totalPages": pagination.TotalPages, "totalItems": pagination.TotalItems, "currentPage": pagination.CurrentPage}

	return c.Status(http.StatusOK).JSON(resp)
}

func (handler *planHandler) CreateHandle(c *fiber.Ctx) (err error) {
	var resp response.Response

	userID, _ := handler.GetUserID(c)
	planDto := new(dto.PlanDTO)

	// collect fields and image file from request as multipart form
	form, formErr := c.MultipartForm()

	if formErr != nil {
		resp.Status = http.StatusExpectationFailed
		resp.Message = formErr.Error()

		return c.Status(http.StatusBadRequest).JSON(resp)
	}

	// parse request body to dto
	if err := c.BodyParser(planDto); err != nil {
		resp.Status = http.StatusExpectationFailed
		resp.Message = err.Error()

		return c.Status(http.StatusBadRequest).JSON(resp)
	}

	fileNames, _ := handler.mediaHelper.Save(c)

	uploadErr := handler.mediaHelper.UploadToAWSS3(fileNames[0])

	if uploadErr != nil {
		resp.Status = http.StatusExpectationFailed
		resp.Message = uploadErr.Error()

		return c.Status(http.StatusBadRequest).JSON(resp)
	}

	planDto.Image = fileNames[0]

	// get other fields from form
	planDto.Title = form.Value["title"][0]
	planDto.Description = form.Value["description"][0]
	planDto.Features = form.Value["features"]

	planDto.CreatedByID = userID
	planDto.Slug = helper.CreateSlug(planDto.Title)

	_, existErr := handler.planService.Read(planDto.Slug)

	if existErr == nil {
		resp.Status = http.StatusConflict
		resp.Message = "Record Already Exist"

		return c.Status(http.StatusConflict).JSON(resp)
	}

	vEs, err := handler.planValidator.Validate(*planDto)

	if err != nil {
		resp.Status = http.StatusUnprocessableEntity
		resp.Message = "validation Error"
		resp.Data = vEs

		return c.Status(http.StatusUnprocessableEntity).JSON(resp)
	}

	if _, err := handler.planService.Create(*planDto); err != nil {
		resp.Status = http.StatusInternalServerError
		resp.Message = err.Error()

		return c.Status(http.StatusInternalServerError).JSON(resp)
	}

	resp.Status = http.StatusCreated
	resp.Message = http.StatusText(http.StatusCreated)
	resp.Data = map[string]interface{}{"plan": planDto}

	return c.Status(http.StatusOK).JSON(resp)
}

func (handler *planHandler) ReadHandle(c *fiber.Ctx) (err error) {
	var resp response.Response

	slug := c.Params("slug")

	planDto, err := handler.planService.Read(slug)

	if err != nil {
		resp.Status = http.StatusNotFound
		resp.Message = "Record Not Found"

		return c.Status(http.StatusNotFound).JSON(resp)
	}

	resp.Status = http.StatusOK
	resp.Message = http.StatusText(http.StatusOK)
	resp.Data = map[string]interface{}{"plan": planDto}

	return c.Status(200).JSON(resp)
}

func (handler *planHandler) UpdateHandle(c *fiber.Ctx) (err error) {
	var resp response.Response
	userID, err := handler.GetUserID(c)

	slug := c.Params("slug")

	if err != nil {
		resp.Status = http.StatusExpectationFailed
		resp.Message = "Exception Error: " + err.Error()

		return c.Status(http.StatusExpectationFailed).JSON(resp)
	}

	planDto, err := handler.planService.Read(slug)
	planDto.Slug = helper.CreateSlug(planDto.Title)

	if err != nil {
		resp.Status = http.StatusNotFound
		resp.Message = "Record Not Found"

		return c.Status(http.StatusNotFound).JSON(resp)
	}

	if err := c.BodyParser(&planDto); err != nil {
		resp.Status = http.StatusExpectationFailed
		resp.Message = err.Error()

		return c.Status(http.StatusExpectationFailed).JSON(resp)
	}

	planDto.Slug = slug
	planDto.UpdatedByID = userID

	vEs, err := handler.planValidator.Validate(planDto)

	if err != nil {
		resp.Status = http.StatusUnprocessableEntity
		resp.Message = "validation Error"
		resp.Data = vEs

		return c.Status(http.StatusUnprocessableEntity).JSON(resp)
	}

	if _, err := handler.planService.Update(planDto); err != nil {
		resp.Status = http.StatusInternalServerError
		resp.Message = err.Error()

		return c.Status(http.StatusInternalServerError).JSON(resp)
	}

	resp.Status = 200
	resp.Message = "Update Successful"

	return c.Status(200).JSON(resp)
}

func (handler *planHandler) DeleteHandle(c *fiber.Ctx) (err error) {

	var resp response.Response
	userID, err := handler.GetUserID(c)

	slug := c.Params("slug")

	if err != nil {
		resp.Status = http.StatusExpectationFailed
		resp.Message = "Exception Error: " + err.Error()

		return c.Status(http.StatusExpectationFailed).JSON(resp)
	}

	_, existsErr := handler.planService.Read(slug)

	if existsErr != nil {
		resp.Status = http.StatusNotFound
		resp.Message = "Record not found"
		return c.Status(http.StatusNotFound).JSON(resp)
	}

	if err := handler.planService.Delete(userID, slug); err != nil {
		resp.Status = http.StatusInternalServerError
		resp.Message = err.Error()

		return c.Status(http.StatusInternalServerError).JSON(resp)
	}

	resp.Status = http.StatusOK
	resp.Message = http.StatusText(http.StatusOK)

	return c.Status(http.StatusOK).JSON(resp)
}
