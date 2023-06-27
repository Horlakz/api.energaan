package app

import (
	"net/http"

	"github.com/gofiber/fiber/v2"

	dto "github.com/horlakz/energaan-api/database/dto/app"
	services "github.com/horlakz/energaan-api/database/services/app"
	"github.com/horlakz/energaan-api/handler"
	contactRequest "github.com/horlakz/energaan-api/payload/request/app/contact"
	"github.com/horlakz/energaan-api/payload/response"
	validators "github.com/horlakz/energaan-api/validator/app"
)

type ContactHandlerInterface interface {
	IndexHandle(c *fiber.Ctx) error
	CreateHandle(c *fiber.Ctx) error
	// UpdateHandle(c *fiber.Ctx) error
	// DeleteHandle(c *fiber.Ctx) error
}

type ContactHandler struct {
	handler.BaseHandler
	contactService   services.ContactServiceInterface
	contactValidator validators.ContactValidator
}

func NewContactHandler(contactService services.ContactServiceInterface) ContactHandlerInterface {
	return &ContactHandler{
		contactService: contactService,
	}
}

func (handler *ContactHandler) IndexHandle(c *fiber.Ctx) (err error) {
	var resp response.Response
	pageable := handler.GeneratePageable(c)
	contacts, pagination, queryError := handler.contactService.ReadAll(pageable)

	if queryError != nil {
		resp.Status = http.StatusUnprocessableEntity
		resp.Message = queryError.Error()

		return c.Status(http.StatusUnprocessableEntity).JSON(resp)
	}

	resp.Status = http.StatusOK
	resp.Message = http.StatusText(http.StatusOK)
	resp.Data = map[string]interface{}{"results": contacts, "totalPages": pagination.TotalPages, "totalItems": pagination.TotalItems, "currentPage": pagination.CurrentPage}

	return c.Status(http.StatusOK).JSON(resp)

}

func (handler *ContactHandler) CreateHandle(c *fiber.Ctx) (err error) {
	var resp response.Response
	var createRequest contactRequest.CreateRequest

	contactDTO := new(dto.ContactDTO)

	// parse request body to dto
	if err := c.BodyParser(&createRequest); err != nil {
		resp.Status = http.StatusExpectationFailed
		resp.Message = err.Error()

		return c.Status(fiber.StatusBadRequest).JSON(resp)
	}

	contactDTO.FullName = createRequest.FullName
	contactDTO.Email = createRequest.Email
	contactDTO.Phone = createRequest.Phone
	contactDTO.Country = createRequest.Country
	contactDTO.Message = createRequest.Message

	vEs, err := handler.contactValidator.Validate(*contactDTO)

	if err != nil {
		resp.Status = http.StatusUnprocessableEntity
		resp.Message = "Invalid Form Data"
		resp.Data = vEs

		return c.Status(http.StatusUnprocessableEntity).JSON(resp)
	}

	_, emailExistsErr := handler.contactService.Read(createRequest.Email)

	if emailExistsErr == nil {
		resp.Status = http.StatusConflict
		resp.Message = "You have already sent a message."

		return c.Status(http.StatusConflict).JSON(resp)
	}

	if _, err := handler.contactService.Create(*contactDTO); err != nil {
		resp.Status = http.StatusInternalServerError
		resp.Message = err.Error()

		return c.Status(http.StatusInternalServerError).JSON(resp)
	}

	resp.Status = http.StatusCreated
	resp.Message = http.StatusText(http.StatusCreated)
	resp.Data = map[string]interface{}{"result": contactDTO}

	return c.Status(http.StatusOK).JSON(resp)
}

// func (handler *ContactHandler) UpdateHandle(c *fiber.Ctx) (err error) {
// 	var resp response.Response
// 	userID, err := handler.GetUserID(c)

// 	id, paramErr := uuid.Parse(c.Params("id"))

// 	if paramErr != nil {
// 		resp.Status = http.StatusExpectationFailed
// 		resp.Message = "Exception Error: " + paramErr.Error()
// 	}

// 	if err != nil {
// 		resp.Status = http.StatusExpectationFailed
// 		resp.Message = "Exception Error: " + err.Error()

// 		return c.Status(http.StatusExpectationFailed).JSON(resp)
// 	}

// 	contactDto, err := handler.contactService.Read(id)

// 	if err != nil {
// 		resp.Status = http.StatusNotFound
// 		resp.Message = "Record Not Found"

// 		return c.Status(http.StatusNotFound).JSON(resp)
// 	}

// 	if err := c.BodyParser(&contactDto); err != nil {
// 		resp.Status = http.StatusExpectationFailed
// 		resp.Message = err.Error()

// 		return c.Status(http.StatusExpectationFailed).JSON(resp)
// 	}

// 	vEs, err := handler.contactValidator.Validate(contactDto)

// 	if err != nil {
// 		resp.Status = http.StatusUnprocessableEntity
// 		resp.Message = "Invalid Form Data"
// 		resp.Data = vEs

// 		return c.Status(http.StatusUnprocessableEntity).JSON(resp)
// 	}

// 	if _, err := handler.contactService.Update(contactDto); err != nil {
// 		resp.Status = http.StatusInternalServerError
// 		resp.Message = err.Error()

// 		return c.Status(http.StatusInternalServerError).JSON(resp)
// 	}

// 	resp.Status = 200
// 	resp.Message = "Update Successful"

// 	return c.Status(200).JSON(resp)
// }

// func (handler *ContactHandler) DeleteHandle(c *fiber.Ctx) (err error) {

// 	var resp response.Response
// 	userID, err := handler.GetUserID(c)

// 	id, paramErr := uuid.Parse(c.Params("id"))

// 	if paramErr != nil {
// 		resp.Status = http.StatusExpectationFailed
// 		resp.Message = "Exception Error: " + paramErr.Error()
// 	}

// 	if err != nil {
// 		resp.Status = http.StatusExpectationFailed
// 		resp.Message = "Exception Error: " + err.Error()

// 		return c.Status(http.StatusExpectationFailed).JSON(resp)
// 	}

// 	_, existsErr := handler.contactService.Read(id)

// 	if existsErr != nil {
// 		resp.Status = http.StatusNotFound
// 		resp.Message = "Record not found"
// 		return c.Status(http.StatusNotFound).JSON(resp)
// 	}

// 	if err := handler.contactService.Delete(userID, id); err != nil {
// 		resp.Status = http.StatusInternalServerError
// 		resp.Message = err.Error()

// 		return c.Status(http.StatusInternalServerError).JSON(resp)
// 	}

// 	resp.Status = http.StatusOK
// 	resp.Message = http.StatusText(http.StatusOK)

// 	return c.Status(http.StatusOK).JSON(resp)
// }
