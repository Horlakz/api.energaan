package router

import (
	"github.com/gofiber/fiber/v2"

	"github.com/horlakz/energaan-api/database"
	repositories "github.com/horlakz/energaan-api/database/repository/app"
	services "github.com/horlakz/energaan-api/database/services/app"
	handlers "github.com/horlakz/energaan-api/handler/app"
	"github.com/horlakz/energaan-api/middleware"
)

func InitializeProductRouter(router fiber.Router, dbConn database.DatabaseInterface) {
	middleware := middleware.Protected()
	productRepository := repositories.NewProductRepository(dbConn)
	categoryRepository := repositories.NewCategoryRepository(dbConn)
	productService := services.NewProductService(productRepository)
	categoryService := services.NewCategoryService(categoryRepository)

	productHandler := handlers.NewProductHandler(productService, categoryService)

	productRoutes := router.Group("/products")

	productRoutes.Get("/", productHandler.IndexHandle)
	productRoutes.Post("/", middleware, productHandler.CreateHandle)
	productRoutes.Get("/:slug", productHandler.ReadHandle)
	productRoutes.Get("/id/:id", productHandler.ReadByUUIDHandle)
	productRoutes.Patch("/:slug", middleware, productHandler.UpdateHandle)
	productRoutes.Delete("/:slug", middleware, productHandler.DeleteHandle)
}
