package auth

import (
	"net/http"

	"github.com/gofiber/fiber/v2"

	dtos "github.com/horlakz/energaan-api/database/dto/auth"
	services "github.com/horlakz/energaan-api/database/services/auth"
	"github.com/horlakz/energaan-api/handler"
	"github.com/horlakz/energaan-api/helper"
	authRequest "github.com/horlakz/energaan-api/payload/request/auth"
	"github.com/horlakz/energaan-api/payload/response"
	authResponse "github.com/horlakz/energaan-api/payload/response/auth"
	"github.com/horlakz/energaan-api/security"
	validators "github.com/horlakz/energaan-api/validator/auth"
)

type AuthHandlerInterface interface {
	LoginHandle(c *fiber.Ctx) error
	RegisterHandle(c *fiber.Ctx) error
}

type authHandler struct {
	handler.BaseHandler
	userService   services.UserServiceInterface
	userValidator validators.UserValidator
	hashUtils     helper.Argon2
}

func NewAuthHandler(userService services.UserServiceInterface) AuthHandlerInterface {
	return &authHandler{
		userService: userService,
	}
}

func (handler *authHandler) LoginHandle(c *fiber.Ctx) error {
	var resp response.Response
	var signinResponse authResponse.LoginResponse
	var signinRequest authRequest.LoginRequest

	if err := c.BodyParser(&signinRequest); err != nil {
		resp.Status = http.StatusExpectationFailed
		resp.Message = "An Error Occured while signin in"
		resp.Data = map[string]interface{}{"Error": err.Error()}

		return c.Status(fiber.StatusBadRequest).JSON(resp)
	}

	user, err := handler.userService.Authenticate(signinRequest.Email, signinRequest.Password)

	if err != nil {
		resp.Status = http.StatusUnprocessableEntity
		resp.Message = err.Error()
		return c.Status(fiber.StatusUnprocessableEntity).JSON(resp)
	}

	accessToken, jwtErr := security.CreateToken(user.UUID.String())

	if jwtErr != nil {
		return jwtErr
	}

	if err != nil {
		return err
	}

	signinResponse.AccessToken = accessToken

	return c.Status(http.StatusOK).JSON(signinResponse)
}

func (handler *authHandler) RegisterHandle(c *fiber.Ctx) error {
	var resp response.Response
	var signupRequest authRequest.RegisterRequest

	if err := c.BodyParser(&signupRequest); err != nil {
		resp.Status = http.StatusExpectationFailed
		resp.Message = "An Error Occured while signing up"
		resp.Data = map[string]interface{}{"Error": err.Error()}

		return c.Status(fiber.StatusBadRequest).JSON(resp)
	}

	var userDto dtos.UserDTO
	userDto.FullName = signupRequest.FullName
	userDto.Email = signupRequest.Email
	userDto.Password, _ = handler.hashUtils.HashPassword(signupRequest.Password)

	vEs, verr := handler.userValidator.Validate(userDto)

	if verr != nil {
		resp.Status = http.StatusUnprocessableEntity
		resp.Message = "Validation Error"
		resp.Data = vEs

		return c.Status(http.StatusUnprocessableEntity).JSON(resp)
	}

	userDto, err := handler.userService.Create(userDto)

	if err != nil {
		resp.Status = http.StatusInternalServerError
		resp.Message = err.Error()

		return c.Status(http.StatusInternalServerError).JSON(resp)
	}

	handler.userService.Create(userDto)

	resp.Status = http.StatusOK
	resp.Message = http.StatusText(http.StatusOK)

	return c.Status(http.StatusOK).JSON(resp)

}
