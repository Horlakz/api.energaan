package main

import (
	"fmt"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"

	"github.com/joho/godotenv"

	"github.com/horlakz/energaan-api/database"
	"github.com/horlakz/energaan-api/router"
)

func main() {
	err := godotenv.Load()

	if err != nil {
		log.Fatal("Error loading .env file")
	} else {
		fmt.Println("Loaded .env file")
	}

	app := fiber.New(fiber.Config{AppName: "Horlakz Email Service v0.0.1"})

	app.Use(recover.New())
	app.Use(logger.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "Origin, Content-Type, Accept",
	}))
	app.Use(limiter.New(limiter.Config{
		Max:               50,
		Expiration:        60 * time.Second,
		LimiterMiddleware: limiter.FixedWindow{},
	}))

	dbConn := database.StartDatabaseClient()

	router.InitializeRouter(app, dbConn)

	log.Fatal(app.Listen(":8000"))
}
