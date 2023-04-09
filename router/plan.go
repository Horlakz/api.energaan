package router

import (
	"github.com/gofiber/fiber/v2"

	"github.com/horlakz/energaan-api/database"
	repositories "github.com/horlakz/energaan-api/database/repository/app"
	services "github.com/horlakz/energaan-api/database/services/app"
	handlers "github.com/horlakz/energaan-api/handler/app"
	"github.com/horlakz/energaan-api/middleware"
)

func InitializePlanRouter(router fiber.Router, dbConn database.DatabaseInterface) {
	middleware := middleware.Protected()
	planRepository := repositories.NewPlanRepository(dbConn)
	planService := services.NewPlanService(planRepository)

	planHandler := handlers.NewPlanHandler(planService)

	planRoutes := router.Group("/plans")

	planRoutes.Get("/", planHandler.IndexHandle)
	planRoutes.Post("/", middleware, planHandler.CreateHandle)
	planRoutes.Get("/:slug", planHandler.ReadHandle)
	planRoutes.Patch("/:slug", middleware, planHandler.UpdateHandle)
	planRoutes.Delete("/:slug", middleware, planHandler.DeleteHandle)
}
