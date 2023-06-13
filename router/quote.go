package router

import (
	"github.com/gofiber/fiber/v2"

	"github.com/horlakz/energaan-api/database"
	repositories "github.com/horlakz/energaan-api/database/repository/app"
	services "github.com/horlakz/energaan-api/database/services/app"
	handlers "github.com/horlakz/energaan-api/handler/app"
	"github.com/horlakz/energaan-api/middleware"
)

func InitializeQuoteRouter(router fiber.Router, dbConn database.DatabaseInterface) {
	middleware := middleware.Protected()

	quoteRepository := repositories.NewQuoteRepository(dbConn)
	productRepository := repositories.NewProductRepository(dbConn)
	planRepository := repositories.NewPlanRepository(dbConn)
	quoteService := services.NewQuoteService(quoteRepository)
	productService := services.NewProductService(productRepository)
	planService := services.NewPlanService(planRepository)

	quoteHandler := handlers.NewQuoteHandler(quoteService, productService, planService)

	quoteRoutes := router.Group("/quotes")

	quoteRoutes.Get("/", middleware, quoteHandler.IndexHandle)
	quoteRoutes.Post("/", quoteHandler.CreateHandle)
	// quoteRoutes.Patch("/:id", middleware, quoteHandler.UpdateHandle)
	// quoteRoutes.Delete("/:id", middleware, quoteHandler.DeleteHandle)
}
