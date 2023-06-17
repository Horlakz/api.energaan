package router

import (
	"github.com/gofiber/fiber/v2"

	"github.com/horlakz/energaan-api/database"
	repositories "github.com/horlakz/energaan-api/database/repository/app"
	userRepository "github.com/horlakz/energaan-api/database/repository/auth"
	services "github.com/horlakz/energaan-api/database/services/app"
	userService "github.com/horlakz/energaan-api/database/services/auth"
	handlers "github.com/horlakz/energaan-api/handler/app"
	"github.com/horlakz/energaan-api/middleware"
)

func InitializeCategoryRouter(router fiber.Router, dbConn database.DatabaseInterface) {
	middleware := middleware.Protected()
	categoryRepository := repositories.NewCategoryRepository(dbConn)
	userRepository := userRepository.NewUserRepository(dbConn)
	categoryService := services.NewCategoryService(categoryRepository)
	userService := userService.NewUserService(userRepository)

	categoryHandler := handlers.NewCategoryHandler(categoryService, userService)

	categoryRoutes := router.Group("/categories")

	categoryRoutes.Get("/", categoryHandler.IndexHandle)
	categoryRoutes.Post("/", middleware, categoryHandler.CreateHandle)
	categoryRoutes.Patch("/:slug", middleware, categoryHandler.UpdateHandle)
	categoryRoutes.Delete("/:slug", middleware, categoryHandler.DeleteHandle)
}
