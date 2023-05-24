package router

import (
	"github.com/gofiber/fiber/v2"

	"github.com/horlakz/energaan-api/database"
	repositories "github.com/horlakz/energaan-api/database/repository/app"
	services "github.com/horlakz/energaan-api/database/services/app"
	handlers "github.com/horlakz/energaan-api/handler/app"
	"github.com/horlakz/energaan-api/middleware"
)

func InitializeCategoryRouter(router fiber.Router, dbConn database.DatabaseInterface) {
	middleware := middleware.Protected()
	categoryRepository := repositories.NewCategoryRepository(dbConn)
	categoryService := services.NewCategoryService(categoryRepository)

	categoryHandler := handlers.NewCategoryHandler(categoryService)

	categoryRoutes := router.Group("/categories")

	categoryRoutes.Get("/", categoryHandler.IndexHandle)
	categoryRoutes.Post("/", middleware, categoryHandler.CreateHandle)
	categoryRoutes.Patch("/:slug", middleware, categoryHandler.UpdateHandle)
	categoryRoutes.Delete("/:slug", middleware, categoryHandler.DeleteHandle)
}
