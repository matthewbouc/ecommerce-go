package api

import (
	"ecommerce/config"
	"ecommerce/internal/api/rest"
	"ecommerce/internal/api/rest/handlers"
	"log"

	"github.com/gofiber/fiber/v3"
)

func StartServer(config config.AppConfig) {
	app := fiber.New()

	router := &rest.Router{
		App: app,
	}

	setupRoutes(router)

	log.Fatal(app.Listen(config.ServerPort))
}

func setupRoutes(router *rest.Router) {
	// user handler
	handlers.SetupUserRoutes(router)
	//catalog handler
	//transactionhandler
}
