package router

import (
	"github.com/gofiber/fiber/v2"

	"github.com/horlakz/energaan-api/database"
	repositories "github.com/horlakz/energaan-api/database/repository/app"
	services "github.com/horlakz/energaan-api/database/services/app"
	handlers "github.com/horlakz/energaan-api/handler/app"
	"github.com/horlakz/energaan-api/middleware"
)

func InitializeContactRouter(router fiber.Router, dbConn database.DatabaseInterface) {
	middleware := middleware.Protected()

	contactRepository := repositories.NewContactRepository(dbConn)
	contactService := services.NewContactService(contactRepository)

	contactHandler := handlers.NewContactHandler(contactService)

	contactRoutes := router.Group("/contacts")

	contactRoutes.Get("/", middleware, contactHandler.IndexHandle)
	contactRoutes.Post("/", contactHandler.CreateHandle)
}
