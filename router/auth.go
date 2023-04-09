package router

import (
	"github.com/gofiber/fiber/v2"

	"github.com/horlakz/energaan-api/database"
	repositories "github.com/horlakz/energaan-api/database/repository/auth"
	services "github.com/horlakz/energaan-api/database/services/auth"
	handlers "github.com/horlakz/energaan-api/handler/auth"
)

func InitializeAuthRouter(router fiber.Router, dbConn database.DatabaseInterface) {

	userRepository := repositories.NewUserRepository(dbConn)
	userService := services.NewUserService(userRepository)

	authHandler := handlers.NewAuthHandler(userService)

	authRoutes := router.Group("/auth")

	authRoutes.Post("/login", authHandler.LoginHandle)
	authRoutes.Post("/register", authHandler.RegisterHandle)
}
