package middleware

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	jwtware "github.com/gofiber/jwt/v2"

	"github.com/horlakz/energaan-api/security"
)

func Protected() fiber.Handler {
	err := jwtware.New(jwtware.Config{
		SigningKey:   []byte(security.GetSecret()),
		ErrorHandler: jwtError,
	})

	if err != nil {
		return err
	}

	return checkUser()
}

func jwtError(c *fiber.Ctx, err error) error {
	if err.Error() == "Missing or malformed JWT" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": err.Error()})
	}

	return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"message": "Invalid or expired JWT"})
}

func checkUser() fiber.Handler {
	return func(c *fiber.Ctx) (err error) {
		_, err = security.ExtractUserID(c.Request())

		if err != nil {
			c.Status(http.StatusForbidden).JSON(fiber.Map{"message": "Unable to extract User"})
		}

		return c.Next()
	}
}
