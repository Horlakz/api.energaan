package app

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	dto "github.com/horlakz/energaan-api/database/dto/app"
	services "github.com/horlakz/energaan-api/database/services/app"
	"github.com/horlakz/energaan-api/handler"
	faqRequest "github.com/horlakz/energaan-api/payload/request/app/faq"
	"github.com/horlakz/energaan-api/payload/response"
	validators "github.com/horlakz/energaan-api/validator/app"
)

type FaqHandlerInterface interface {
	IndexHandle(c *fiber.Ctx) error
	CreateHandle(c *fiber.Ctx) error
	UpdateHandle(c *fiber.Ctx) error
	DeleteHandle(c *fiber.Ctx) error
}

type FaqHandler struct {
	handler.BaseHandler
	faqService   services.FaqServiceInterface
	faqValidator validators.FaqValidator
}

func NewFaqHandler(faqService services.FaqServiceInterface) FaqHandlerInterface {
	return &FaqHandler{
		faqService: faqService,
	}
}

func (handler *FaqHandler) IndexHandle(c *fiber.Ctx) (err error) {
	var resp response.Response
	faqs, queryError := handler.faqService.ReadAll()

	if queryError != nil {
		resp.Status = http.StatusUnprocessableEntity
		resp.Message = queryError.Error()

		return c.Status(http.StatusUnprocessableEntity).JSON(resp)
	}

	resp.Status = http.StatusOK
	resp.Message = http.StatusText(http.StatusOK)

	if len(faqs) < 1 {
		resp.Data = map[string]interface{}{"result": []string{}}
		return c.Status(http.StatusOK).JSON(resp)
	}

	resp.Data = map[string]interface{}{"result": faqs}
	return c.Status(http.StatusOK).JSON(resp)

}

func (handler *FaqHandler) CreateHandle(c *fiber.Ctx) (err error) {
	var resp response.Response
	var createRequest faqRequest.CreateRequest

	userID, _ := handler.GetUserID(c)
	faqDTO := new(dto.FaqDTO)

	// parse request body to dto
	if err := c.BodyParser(&createRequest); err != nil {
		resp.Status = http.StatusExpectationFailed
		resp.Message = err.Error()

		return c.Status(fiber.StatusBadRequest).JSON(resp)
	}

	faqDTO.Title = createRequest.Title
	faqDTO.Description = createRequest.Description
	faqDTO.CreatedByID = userID

	vEs, err := handler.faqValidator.Validate(*faqDTO)

	if err != nil {
		resp.Status = http.StatusUnprocessableEntity
		resp.Message = "validation Error"
		resp.Data = vEs

		return c.Status(http.StatusUnprocessableEntity).JSON(resp)
	}

	if _, err := handler.faqService.Create(*faqDTO); err != nil {
		resp.Status = http.StatusInternalServerError
		resp.Message = err.Error()

		return c.Status(http.StatusInternalServerError).JSON(resp)
	}

	resp.Status = http.StatusCreated
	resp.Message = http.StatusText(http.StatusCreated)
	resp.Data = map[string]interface{}{"result": faqDTO}

	return c.Status(http.StatusOK).JSON(resp)
}

func (handler *FaqHandler) UpdateHandle(c *fiber.Ctx) (err error) {
	var resp response.Response
	userID, err := handler.GetUserID(c)

	id, paramErr := uuid.Parse(c.Params("id"))

	if paramErr != nil {
		resp.Status = http.StatusExpectationFailed
		resp.Message = "Exception Error: " + paramErr.Error()
	}

	if err != nil {
		resp.Status = http.StatusExpectationFailed
		resp.Message = "Exception Error: " + err.Error()

		return c.Status(http.StatusExpectationFailed).JSON(resp)
	}

	faqDto, err := handler.faqService.Read(id)

	if err != nil {
		resp.Status = http.StatusNotFound
		resp.Message = "Record Not Found"

		return c.Status(http.StatusNotFound).JSON(resp)
	}

	if err := c.BodyParser(&faqDto); err != nil {
		resp.Status = http.StatusExpectationFailed
		resp.Message = err.Error()

		return c.Status(http.StatusExpectationFailed).JSON(resp)
	}

	faqDto.UpdatedByID = userID

	vEs, err := handler.faqValidator.Validate(faqDto)

	if err != nil {
		resp.Status = http.StatusUnprocessableEntity
		resp.Message = "validation Error"
		resp.Data = vEs

		return c.Status(http.StatusUnprocessableEntity).JSON(resp)
	}

	if _, err := handler.faqService.Update(faqDto); err != nil {
		resp.Status = http.StatusInternalServerError
		resp.Message = err.Error()

		return c.Status(http.StatusInternalServerError).JSON(resp)
	}

	resp.Status = 200
	resp.Message = "Update Successful"

	return c.Status(200).JSON(resp)
}

func (handler *FaqHandler) DeleteHandle(c *fiber.Ctx) (err error) {

	var resp response.Response
	userID, err := handler.GetUserID(c)

	id, paramErr := uuid.Parse(c.Params("id"))

	if paramErr != nil {
		resp.Status = http.StatusExpectationFailed
		resp.Message = "Exception Error: " + paramErr.Error()
	}

	if err != nil {
		resp.Status = http.StatusExpectationFailed
		resp.Message = "Exception Error: " + err.Error()

		return c.Status(http.StatusExpectationFailed).JSON(resp)
	}

	_, existsErr := handler.faqService.Read(id)

	if existsErr != nil {
		resp.Status = http.StatusNotFound
		resp.Message = "Record not found"
		return c.Status(http.StatusNotFound).JSON(resp)
	}

	if err := handler.faqService.Delete(userID, id); err != nil {
		resp.Status = http.StatusInternalServerError
		resp.Message = err.Error()

		return c.Status(http.StatusInternalServerError).JSON(resp)
	}

	resp.Status = http.StatusOK
	resp.Message = http.StatusText(http.StatusOK)

	return c.Status(http.StatusOK).JSON(resp)
}
