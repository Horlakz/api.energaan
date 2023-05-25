package app

import (
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

type GalleryHandlerInterface interface {
	IndexHandle(c *fiber.Ctx) error
	CreateHandle(c *fiber.Ctx) error
	UpdateHandle(c *fiber.Ctx) error
	DeleteHandle(c *fiber.Ctx) error
}

type galleryHandler struct {
	handler.BaseHandler
	mediaHelper      helper.Media
	galleryService   services.GalleryServiceInterface
	galleryValidator validators.GalleryValidator
}

func NewgalleryHandler(galleryService services.GalleryServiceInterface) GalleryHandlerInterface {
	return &galleryHandler{
		galleryService: galleryService,
	}
}

func (handler *galleryHandler) IndexHandle(c *fiber.Ctx) (err error) {
	var resp response.Response
	categories, queryError := handler.galleryService.ReadAll()

	if queryError != nil {
		resp.Status = http.StatusUnprocessableEntity
		resp.Message = queryError.Error()

		return c.Status(http.StatusUnprocessableEntity).JSON(resp)
	}

	resp.Status = http.StatusOK
	resp.Message = http.StatusText(http.StatusOK)

	if len(categories) < 1 {
		resp.Data = map[string]interface{}{"result": []string{}}
		return c.Status(http.StatusOK).JSON(resp)
	}

	resp.Data = map[string]interface{}{"result": categories}
	return c.Status(http.StatusOK).JSON(resp)

}

func (handler *galleryHandler) CreateHandle(c *fiber.Ctx) (err error) {
	var resp response.Response

	userID, _ := handler.GetUserID(c)
	galleryDTO := new(dto.GalleryDTO)

	// collect fields and image file from request as multipart form
	form, formErr := c.MultipartForm()

	if formErr != nil {
		resp.Status = http.StatusExpectationFailed
		resp.Message = formErr.Error()

		return c.Status(http.StatusBadRequest).JSON(resp)
	}

	// parse request body to dto
	if err := c.BodyParser(galleryDTO); err != nil {
		resp.Status = http.StatusExpectationFailed
		resp.Message = err.Error()

		return c.Status(http.StatusBadRequest).JSON(resp)
	}

	fileNames, _ := handler.mediaHelper.Save(c)

	galleryDTO.Image = fileNames[0]
	galleryDTO.Title = form.Value["title"][0]
	galleryDTO.CreatedByID = userID

	vEs, err := handler.galleryValidator.Validate(*galleryDTO)

	if err != nil {
		resp.Status = http.StatusUnprocessableEntity
		resp.Message = "validation Error"
		resp.Data = vEs

		return c.Status(http.StatusUnprocessableEntity).JSON(resp)
	}

	if _, err := handler.galleryService.Create(*galleryDTO); err != nil {
		resp.Status = http.StatusInternalServerError
		resp.Message = err.Error()

		return c.Status(http.StatusInternalServerError).JSON(resp)
	}

	resp.Status = http.StatusCreated
	resp.Message = http.StatusText(http.StatusCreated)
	resp.Data = map[string]interface{}{"result": galleryDTO}

	return c.Status(http.StatusOK).JSON(resp)
}

func (handler *galleryHandler) UpdateHandle(c *fiber.Ctx) (err error) {
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

	galleryDto, err := handler.galleryService.Read(id)

	if err != nil {
		resp.Status = http.StatusNotFound
		resp.Message = "Record Not Found"

		return c.Status(http.StatusNotFound).JSON(resp)
	}

	if err := c.BodyParser(&galleryDto); err != nil {
		resp.Status = http.StatusExpectationFailed
		resp.Message = err.Error()

		return c.Status(http.StatusExpectationFailed).JSON(resp)
	}

	galleryDto.UpdatedByID = userID

	vEs, err := handler.galleryValidator.Validate(galleryDto)

	if err != nil {
		resp.Status = http.StatusUnprocessableEntity
		resp.Message = "validation Error"
		resp.Data = vEs

		return c.Status(http.StatusUnprocessableEntity).JSON(resp)
	}

	if _, err := handler.galleryService.Update(galleryDto); err != nil {
		resp.Status = http.StatusInternalServerError
		resp.Message = err.Error()

		return c.Status(http.StatusInternalServerError).JSON(resp)
	}

	resp.Status = 200
	resp.Message = "Update Successful"

	return c.Status(200).JSON(resp)
}

func (handler *galleryHandler) DeleteHandle(c *fiber.Ctx) (err error) {

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

	_, existsErr := handler.galleryService.Read(id)

	if existsErr != nil {
		resp.Status = http.StatusNotFound
		resp.Message = "Record not found"
		return c.Status(http.StatusNotFound).JSON(resp)
	}

	if err := handler.galleryService.Delete(userID, id); err != nil {
		resp.Status = http.StatusInternalServerError
		resp.Message = err.Error()

		return c.Status(http.StatusInternalServerError).JSON(resp)
	}

	resp.Status = http.StatusOK
	resp.Message = http.StatusText(http.StatusOK)

	return c.Status(http.StatusOK).JSON(resp)
}
