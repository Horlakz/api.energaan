package router

import (
	"github.com/gofiber/fiber/v2"

	"github.com/horlakz/energaan-api/database"
	repositories "github.com/horlakz/energaan-api/database/repository/app"
	services "github.com/horlakz/energaan-api/database/services/app"
	handlers "github.com/horlakz/energaan-api/handler/app"
	"github.com/horlakz/energaan-api/middleware"
)

func InitializeFaqRouter(router fiber.Router, dbConn database.DatabaseInterface) {
	middleware := middleware.Protected()
	faqRepository := repositories.NewFaqRepository(dbConn)
	faqService := services.NewFaqService(faqRepository)

	faqHandler := handlers.NewFaqHandler(faqService)

	faqRoutes := router.Group("/faqs")

	faqRoutes.Get("/", faqHandler.IndexHandle)
	faqRoutes.Post("/", middleware, faqHandler.CreateHandle)
	faqRoutes.Patch("/:id", middleware, faqHandler.UpdateHandle)
	faqRoutes.Delete("/:id", middleware, faqHandler.DeleteHandle)
}
