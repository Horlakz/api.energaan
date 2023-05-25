package router

import (
	"github.com/gofiber/fiber/v2"

	"github.com/horlakz/energaan-api/database"
	repositories "github.com/horlakz/energaan-api/database/repository/app"
	services "github.com/horlakz/energaan-api/database/services/app"
	handlers "github.com/horlakz/energaan-api/handler/app"
	"github.com/horlakz/energaan-api/middleware"
)

func InitializeGalleryRouter(router fiber.Router, dbConn database.DatabaseInterface) {
	middleware := middleware.Protected()
	galleryRepository := repositories.NewGalleryRepository(dbConn)
	galleryService := services.NewGalleryService(galleryRepository)

	galleryHandler := handlers.NewgalleryHandler(galleryService)

	galleryRoutes := router.Group("/gallery")

	galleryRoutes.Get("/", galleryHandler.IndexHandle)
	galleryRoutes.Post("/", middleware, galleryHandler.CreateHandle)
	galleryRoutes.Patch("/:id", middleware, galleryHandler.UpdateHandle)
	galleryRoutes.Delete("/:id", middleware, galleryHandler.DeleteHandle)
}
