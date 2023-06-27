package app

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	dto "github.com/horlakz/energaan-api/database/dto/app"
	services "github.com/horlakz/energaan-api/database/services/app"
	"github.com/horlakz/energaan-api/handler"
	quoteRequest "github.com/horlakz/energaan-api/payload/request/app/quote"
	"github.com/horlakz/energaan-api/payload/response"
	validators "github.com/horlakz/energaan-api/validator/app"
)

type QuoteResult struct {
	Quote   dto.QuoteDTO `json:"quote"`
	Service interface{}  `json:"serviceDetails"`
}

type QuoteHandlerInterface interface {
	IndexHandle(c *fiber.Ctx) error
	CreateHandle(c *fiber.Ctx) error
	// UpdateHandle(c *fiber.Ctx) error
	// DeleteHandle(c *fiber.Ctx) error
}

type QuoteHandler struct {
	handler.BaseHandler
	quoteService   services.QuoteServiceInterface
	productService services.ProductServiceInterface
	planService    services.PlanServiceInterface
	quoteValidator validators.QuoteValidator
}

func NewQuoteHandler(
	quoteService services.QuoteServiceInterface,
	productService services.ProductServiceInterface,
	planService services.PlanServiceInterface,
) QuoteHandlerInterface {
	return &QuoteHandler{
		quoteService:   quoteService,
		productService: productService,
		planService:    planService,
	}
}

func (handler *QuoteHandler) IndexHandle(c *fiber.Ctx) (err error) {
	var resp response.Response
	pageable := handler.GeneratePageable(c)
	quotes, pagination, queryError := handler.quoteService.ReadAll(pageable)

	if queryError != nil {
		resp.Status = http.StatusUnprocessableEntity
		resp.Message = queryError.Error()

		return c.Status(http.StatusUnprocessableEntity).JSON(resp)
	}

	results := make([]QuoteResult, len(quotes))

	// create a new object of response, loop through the quotes, check the serviceType, if type is product, call productService.ReadByUUID, if type is plan, call planService.ReadByUUID and attach the result to the quote object
	for i, quote := range quotes {
		if quote.ServiceType == "product" {
			product, err := handler.productService.ReadByUUID(quote.ServiceId)

			if err != nil {
				results[i].Quote = quote
				results[i].Service = nil
				continue
			}

			results[i].Quote = quote
			results[i].Service = product
		} else if quote.ServiceType == "plan" {
			plan, err := handler.planService.ReadByUUID(quote.ServiceId)

			if err != nil {
				results[i].Quote = quote
				results[i].Service = nil
				continue
			}

			results[i].Quote = quote
			results[i].Service = plan
		}
	}

	resp.Status = http.StatusOK
	resp.Message = http.StatusText(http.StatusOK)
	resp.Data = map[string]interface{}{
		"result":      results,
		"totalPages":  pagination.TotalPages,
		"totalItems":  pagination.TotalItems,
		"currentPage": pagination.CurrentPage,
	}

	return c.Status(http.StatusOK).JSON(resp)

}

func (handler *QuoteHandler) CreateHandle(c *fiber.Ctx) (err error) {
	var resp response.Response
	var createRequest quoteRequest.CreateRequest

	quoteDTO := new(dto.QuoteDTO)

	// parse request body to dto
	if err := c.BodyParser(&createRequest); err != nil {
		resp.Status = http.StatusExpectationFailed
		resp.Message = err.Error()

		return c.Status(fiber.StatusBadRequest).JSON(resp)
	}

	serviceId, err := uuid.Parse(createRequest.ServiceId)

	if err != nil {
		resp.Status = http.StatusExpectationFailed
		resp.Message = err.Error()

		return c.Status(http.StatusExpectationFailed).JSON(resp)
	}

	quoteDTO.FullName = createRequest.FullName
	quoteDTO.Email = createRequest.Email
	quoteDTO.Phone = createRequest.Phone
	quoteDTO.Country = createRequest.Country
	quoteDTO.ServiceId = serviceId
	quoteDTO.ServiceType = createRequest.ServiceType

	vEs, err := handler.quoteValidator.Validate(*quoteDTO)

	if err != nil {
		resp.Status = http.StatusUnprocessableEntity
		resp.Message = "Invalid Form Data"
		resp.Data = vEs

		return c.Status(http.StatusUnprocessableEntity).JSON(resp)
	}

	alreadyExists, _ := handler.quoteService.CheckEmailServiceTypeExists(createRequest.Email, createRequest.ServiceType)

	if alreadyExists {
		resp.Status = http.StatusConflict
		resp.Message = "You have already requested for this service"

		return c.Status(http.StatusConflict).JSON(resp)
	}

	if _, err := handler.quoteService.Create(*quoteDTO); err != nil {
		resp.Status = http.StatusInternalServerError
		resp.Message = err.Error()

		return c.Status(http.StatusInternalServerError).JSON(resp)
	}

	resp.Status = http.StatusCreated
	resp.Message = http.StatusText(http.StatusCreated)
	resp.Data = map[string]interface{}{"result": quoteDTO}

	return c.Status(http.StatusOK).JSON(resp)
}

// func (handler *QuoteHandler) UpdateHandle(c *fiber.Ctx) (err error) {
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

// 	quoteDto, err := handler.quoteService.Read(id)

// 	if err != nil {
// 		resp.Status = http.StatusNotFound
// 		resp.Message = "Record Not Found"

// 		return c.Status(http.StatusNotFound).JSON(resp)
// 	}

// 	if err := c.BodyParser(&quoteDto); err != nil {
// 		resp.Status = http.StatusExpectationFailed
// 		resp.Message = err.Error()

// 		return c.Status(http.StatusExpectationFailed).JSON(resp)
// 	}

// 	vEs, err := handler.quoteValidator.Validate(quoteDto)

// 	if err != nil {
// 		resp.Status = http.StatusUnprocessableEntity
// 		resp.Message = "Invalid Form Data"
// 		resp.Data = vEs

// 		return c.Status(http.StatusUnprocessableEntity).JSON(resp)
// 	}

// 	if _, err := handler.quoteService.Update(quoteDto); err != nil {
// 		resp.Status = http.StatusInternalServerError
// 		resp.Message = err.Error()

// 		return c.Status(http.StatusInternalServerError).JSON(resp)
// 	}

// 	resp.Status = 200
// 	resp.Message = "Update Successful"

// 	return c.Status(200).JSON(resp)
// }

// func (handler *QuoteHandler) DeleteHandle(c *fiber.Ctx) (err error) {

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

// 	_, existsErr := handler.quoteService.Read(id)

// 	if existsErr != nil {
// 		resp.Status = http.StatusNotFound
// 		resp.Message = "Record not found"
// 		return c.Status(http.StatusNotFound).JSON(resp)
// 	}

// 	if err := handler.quoteService.Delete(userID, id); err != nil {
// 		resp.Status = http.StatusInternalServerError
// 		resp.Message = err.Error()

// 		return c.Status(http.StatusInternalServerError).JSON(resp)
// 	}

// 	resp.Status = http.StatusOK
// 	resp.Message = http.StatusText(http.StatusOK)

// 	return c.Status(http.StatusOK).JSON(resp)
// }
